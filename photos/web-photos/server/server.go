package server

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/graeme-hill/gnet/photos/events"
	"github.com/graeme-hill/gnet/sys/filestore"
	"github.com/graeme-hill/gnet/sys/gnet"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/graeme-hill/gnet/sys/pb"
)

type Options struct {
	Addr                string
	DomainEventsRPCAddr string
	EventStoreConnStr   string
	FileStoreConnStr    string
}

func makeFileName(hash string) string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return "ph_" + id.String() + "_" + hash
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func handleUpload(opts Options, w http.ResponseWriter, r *http.Request) {
	filesDB := filestore.NewFileStoreConn(opts.EventStoreConnStr)

	// Validate
	if r.Method != http.MethodPost {
		http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	hash, hasHash := vars["name"]
	if !hasHash {
		http.Error(w, "missing file hash", http.StatusBadRequest)
		return
	}

	// Write the file somewhere first
	fileName := makeFileName(hash)
	err := filesDB.Write(fileName, r.Body)

	if err != nil {
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	// Record domain event
	de, err := events.NewPhotoUploadedEvent(hash)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	conn, err := grpc.Dial(opts.DomainEventsRPCAddr, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	c := pb.NewDomainEventsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = c.InsertDomainEvent(ctx, &pb.InsertDomainEventRequest{
		Type: de.Type,
		Data: []byte{1, 2},
	})

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handler(opts Options, inner func(Options, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inner(opts, w, r)
	}
}

func formatAddr(addr string) string {
	parts := strings.Split(addr, ":")
	return ":" + parts[len(parts)-1] + "::"
}

func waitUntilPingable(ctx context.Context, addr string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "ping")
	pingURL := u.String()

	for i := 0; i < 20; i++ {
		response, err := http.Get(pingURL)
		if err == nil && response.StatusCode == http.StatusOK {
			return nil
		}

		select {
		case <-ctx.Done():
			return errors.New("context canceled before server became pingable")
		case <-time.After(100 * time.Millisecond):
		}
	}

	return errors.New("gave up waiting for photos web api to come online")
}

func Run(ctx context.Context, opts Options) gnet.Service {
	router := mux.NewRouter()
	router.HandleFunc("/upload/{hash}", handler(opts, handleUpload))
	router.HandleFunc("/ping", handlePing)

	srv := &http.Server{
		Handler:      router,
		Addr:         formatAddr(opts.Addr),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	errChan := make(chan error)
	startedChan := make(chan struct{})

	// run the server
	go func() {
		err := srv.ListenAndServe()
		if err == http.ErrServerClosed {
			// ErrServerClosed isn't an error!
			err = nil
		}
		errChan <- err
		log.Printf("OFFLINE - photos API: %s", opts.Addr)
	}()

	// tell the server to stop when canceled
	go func() {
		select {
		case <-ctx.Done():
			_ = srv.Close()
		}
	}()

	// detect when actually running
	go func() {
		err := waitUntilPingable(ctx, opts.Addr)
		if err != nil {
			errChan <- errors.Wrap(err, "error when waiting for server to become pingable")
			return
		}
		log.Printf("ONLINE - photos API: %s", opts.Addr)
		startedChan <- struct{}{}
	}()

	return gnet.Service{
		Over:    errChan,
		Running: startedChan,
	}
}
