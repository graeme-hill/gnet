package filestore

import "io"

// FileStore is S3-like file access interface.
type FileStore interface {
	Read(path string) (io.Reader, bool, error)
	Write(path string, data io.Reader) error
	Delete(path string) (bool, error)
}
