package keyvaluestore

var NewKeyValueStoreConn func(string) KeyValueStore = NewInMemKeyValueStore

type KeyValueStore interface {
	Get(string) ([]byte, bool, error)
	Set(string, []byte) error
	Delete(string) error
	DeleteMany(string) error
	ReadMany(min string, max string, limit int) (Cursor, error)
}

type Cursor interface {
	Next() (ReadResult, error)
}

type ReadResult struct {
	Key   string
	Value []byte
	More  bool
}
