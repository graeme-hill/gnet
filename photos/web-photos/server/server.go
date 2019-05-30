package server

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/graeme-hill/gnet/photos/events"
	"github.com/graeme-hill/gnet/sys/filestore"
	"github.com/oklog/ulid"
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

func Run(ctx context.Context, opts Options) <-chan error {
	router := mux.NewRouter()
	router.HandleFunc("/upload/{hash}", handler(opts, handleUpload))

	srv := &http.Server{
		Handler:      router,
		Addr:         opts.Addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	errChan := make(chan error)

	// run the server
	go func() {
		log.Println("WEB SERVER: running listenandserve " + opts.Addr)
		errChan <- srv.ListenAndServe()
		log.Println("WEB SERVER: shutted down")
	}()

	// tell the server to stop when canceled
	go func() {
		log.Println("WEB SERVER: waiting for ctx to be done")
		select {
		case <-ctx.Done():
			log.Println("WEB SERVER: ctx is done")
			_ = srv.Close()
		}
	}()

	return errChan
}
