package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/graeme-hill/gnet/photos/events"
	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/filestore"
	"github.com/oklog/ulid"

	depb "github.com/graeme-hill/sys/rpc-domainevents/pb"
)

var filesDB = filestore.NewFileStoreConn(":memory:")
var eventsDB = eventstore.NewEventStoreConn(":memory:")

func makeFileName(hash string) string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return "ph_" + id.String() + "_" + hash
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/upload/{hash}", func(w http.ResponseWriter, r *http.Request) {
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
		c := depb.NewDomainEventsClient(scanClient.conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = c.InsertDomainEvent(ctx, &depb.InsertDomainEventRequest{
			Type: de.Type,
			Data: []byte{1,2},
		})

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
