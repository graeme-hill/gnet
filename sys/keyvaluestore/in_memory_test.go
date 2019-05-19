package keyvaluestore

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestInMemKeyValStore(t *testing.T) {
	kvs := NewInMemKeyValueStore("mem")

	data0 := []byte{1, 2, 3}
	data1 := []byte{2, 4, 6}
	data2 := []byte{3, 6, 9}

	// set
	err := kvs.Set("foo:a", data0)
	require.NoError(t, err)
	err = kvs.Set("foo:b", data1)
	require.NoError(t, err)
	err = kvs.Set("foo:c", data2)
	require.NoError(t, err)

	// get
	bytes, ok, err := kvs.Get("foo:a")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, data0, bytes)

	bytes, ok, err = kvs.Get("foo:b")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, data1, bytes)

	bytes, ok, err = kvs.Get("notfound")
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, bytes)

	// readmany high limit
	reader, err := kvs.ReadFrom("foo:", 10)
	require.NoError(t, err)
	found := map[string][]byte{}
	for reader.Next() {
		k, err := reader.Key()
		require.NoError(t, err)
		v, err := reader.Value()
		require.NoError(t, err)
		found[k] = v
	}
	require.Equal(t, map[string][]byte{
		"foo:a": data0,
		"foo:b": data1,
		"foo:c": data2,
	}, found)

	// readmany low limit
	reader, err = kvs.ReadFrom("foo:a", 1)
	require.NoError(t, err)
	found = map[string][]byte{}
	for reader.Next() {
		k, err := reader.Key()
		require.NoError(t, err)
		v, err := reader.Value()
		require.NoError(t, err)
		found[k] = v
	}
	require.Equal(t, map[string][]byte{
		"foo:a": data0,
	}, found)

	// readrange
	reader, err = kvs.ReadRange("foo:a", "foo:b", 10)
	require.NoError(t, err)
	found = map[string][]byte{}
	for reader.Next() {
		k, err := reader.Key()
		require.NoError(t, err)
		v, err := reader.Value()
		require.NoError(t, err)
		found[k] = v
	}
	require.Equal(t, map[string][]byte{
		"foo:a": data0,
		"foo:b": data1,
	}, found)
}