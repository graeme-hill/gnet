package filestore

var NewFileStoreConn func(string) FileStore = NewInMemFileStore
