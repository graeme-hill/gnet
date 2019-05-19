package keyvaluestore

var NewKeyValueStoreConn func(string) KeyValueStore = NewInMemKeyValueStore

type KeyValueStore interface {
	Get(string) ([]byte, bool, error)
	Set(string, []byte) error
	Delete(string) error
	DeleteMany(string) error
	ReadRange(min, max string, limit int) (Cursor, error)
	ReadFrom(start string, limit int) (Cursor, error)
}

type Cursor interface {
	Next() bool
	Key() (string, error)
	Value() ([]byte, error)
}
