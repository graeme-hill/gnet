package keyvaluestore

import (
	"sort"
	"strings"
)

func NewInMemKeyValueStore(connStr string) KeyValueStore {
	return &inMemKeyValueStore{
		data: map[string][]byte{},
	}
}

type inMemKeyValueStore struct {
	data map[string][]byte
}

func (kvs *inMemKeyValueStore) Get(key string) ([]byte, bool, error) {
	value, hasValue := kvs.data[key]
	if !hasValue {
		return nil, false, nil
	}
	return value, true, nil
}

func (kvs *inMemKeyValueStore) Set(key string, data []byte) error {
	kvs.data[key] = data
	return nil
}

func (kvs *inMemKeyValueStore) Delete(key string) error {
	delete(kvs.data, key)
	return nil
}

func (kvs *inMemKeyValueStore) DeleteMany(prefix string) error {
	for key := range kvs.data {
		if strings.HasPrefix(key, prefix) {
			kvs.Delete(key)
		}
	}
	return nil
}

func (kvs *inMemKeyValueStore) ReadMany(min string, max string, limit int) (Cursor, error) {
	keys := []string{}
	for k := range kvs.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := []keyValuePair{}

	for _, k := range keys {
		if strings.Compare(k, min) >= 0 && strings.Compare(k, max) <= 0 {
			pairs = append(pairs, keyValuePair{
				key:   k,
				value: kvs.data[k],
			})
		}

		if limit > 0 && len(pairs) >= limit {
			break
		}
	}

	return &inMemCursor{
		pairs:   pairs,
		current: 0,
	}, nil
}

type keyValuePair struct {
	key   string
	value []byte
}

type inMemCursor struct {
	pairs   []keyValuePair
	current int
}

func (c *inMemCursor) Next() (ReadResult, error) {
	pair := c.pairs[c.current]
	c.current += 1
	return ReadResult{
		Key:   pair.key,
		Value: pair.value,
		More:  len(c.pairs) >= c.current,
	}, nil
}
