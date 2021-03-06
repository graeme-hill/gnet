package filestore

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

var stores = map[string]map[string][]byte{}

// InMemFileStore just stores bytes in a map.
type inMemFileStore struct {
	files map[string][]byte
}

// Read opens a reader on the file if it exists.
func (fs *inMemFileStore) Read(path string) (io.Reader, bool, error) {
	fileBytes, ok := fs.files[path]
	if !ok {
		return nil, false, nil
	}

	return bytes.NewReader(fileBytes), true, nil
}

// Write adds file to map.
func (fs *inMemFileStore) Write(path string, data io.Reader) error {
	fileBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return errors.Wrap(err, "Failed to read all bytes")
	}
	fs.files[path] = fileBytes
	return nil
}

func (fs *inMemFileStore) Delete(path string) (bool, error) {
	_, exists := fs.files[path]
	if !exists {
		return false, nil
	}
	delete(fs.files, path)
	return true, nil
}

// NewInMemFileStore uses a map to store bytes. All data is erased when process
// dies. Multiple calls with same connection string reference same data.
func NewInMemFileStore(connStr string) FileStore {
	store, exists := stores[connStr]
	if !exists {
		store = map[string][]byte{}
	}
	return &inMemFileStore{
		files: store,
	}
}
