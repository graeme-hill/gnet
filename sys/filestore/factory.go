package filestore

func NewFileStoreConn(connStr string) FileStore {
	// when there is a non in-memory implementation make this smarter
	return NewInMemFileStore(connStr)
}
