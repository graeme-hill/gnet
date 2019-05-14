package filestore

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// InMemFileStore just stores bytes in a map.
type InMemFileStore struct {
	files map[string][]byte
}

// Read opens a reader on the file if it exists.
func (fs *InMemFileStore) Read(path string) (io.Reader, error) {
	fileBytes, ok := fs.files[path]
	if !ok {
		return nil, errors.New("file not found")
	}

	return bytes.NewReader(fileBytes), nil
}

// Write adds file to map.
func (fs *InMemFileStore) Write(path string, data io.Reader) error {
	fileBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return errors.Wrap(err, "Failed to read all bytes")
	}
	fs.files[path] = fileBytes
	return nil
}

func NewInMemFileStore() FileStore {
	return &InMemFileStore{
		files: map[string][]byte{},
	}
}
