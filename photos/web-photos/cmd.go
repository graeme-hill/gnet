package main

import (
	"net/http"
	"github.com/graeme-hill/gnet/lib/filestore"
	"github.com/graeme-hill/gnet/lib/eventstore"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid"
)

var store = filestore.NewInMemFileStore()
var events = eventstore.NewInMemFileStore()

const MAX_FILE_NAME_LENGTH = 100

func makeFileName(original string) string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	remaining_chars = MAX_FILE_NAME_LENGTH - len(id) - 1
	return id + "_" + original[:remaining_chars]
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/upload/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Validate
		if r.Method != http.MethodPost {
			http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		name, hasName := vars["name"]
		if !hasName {
			http.Error(w, "missing file name", http.StatusBadRequest)
			return
		}

		// Write the file somewhere first
		fileName := makeFileName(name)
		err := store.Write(name, r.Body)

		if err != nil {
			http.Error(w, "failed to upload file", http.StatusInternalServerError)
			return
		}

		// Record domain event
		events.Insert(DomainEvent{
			Type: "upload_photo",
			Data: 
		})

		w.WriteHeader(http.StatusOK)
	})
}
