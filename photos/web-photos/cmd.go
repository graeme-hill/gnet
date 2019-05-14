package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/filestore"
	"github.com/oklog/ulid"
)

var files = filestore.NewFileStoreConn("")
var events = eventstore.NewEventStoreConn()

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
		err := files.Write(fileName, r.Body)

		if err != nil {
			http.Error(w, "failed to upload file", http.StatusInternalServerError)
			return
		}

		// Record domain event
		//events.Insert(...)

		w.WriteHeader(http.StatusOK)
	})
}
