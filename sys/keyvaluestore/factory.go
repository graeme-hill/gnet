package keyvaluestore

func NewKeyValueStore(connStr string) KeyValueStore {
	// when there is a non im-memory implementation make this smarter
	return NewInMemKeyValueStore(connStr)
}