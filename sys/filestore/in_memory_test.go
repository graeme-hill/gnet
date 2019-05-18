package filestore

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	fs := NewInMemFileStore("mem")

	data1 := []byte{1, 2, 3}
	data2 := []byte{2, 4, 6, 8}

	fs.Write("/foo/bar", bytes.NewReader(data1))
	fs.Write("/foo/baroo", bytes.NewReader(data2))

	reader, found, err := fs.Read("/foo/bar")
	require.NoError(t, err)
	require.True(t, found)
	bytes, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, data1, bytes)

	reader, found, err = fs.Read("/foo/baroo")
	require.NoError(t, err)
	require.True(t, found)
	bytes, err = ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, data2, bytes)

	reader, found, err = fs.Read("/foo/nahh")
	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, reader)

	deleted, err := fs.Delete("/foo/nahh")
	require.NoError(t, err)
	require.False(t, deleted)

	deleted, err = fs.Delete("/foo/bar")
	require.NoError(t, err)
	require.True(t, deleted)

	reader, found, err = fs.Read("/foo/bar")
	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, reader)
}
