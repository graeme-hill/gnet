package filestore

var NewFileStoreConn func() FileStore = NewInMemFileStore
