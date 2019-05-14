package eventstore

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type NewUserEvent struct {
	name string
}

func (e NewUserEvent) Type() string {
	return "new_user"
}

func TestInMem(t *testing.T) {
	store := NewInMemEventStore()
	graeme := []byte{1}
	foobar := []byte{2}
	store.Insert(DomainEvent{Type: "new_user", Data: graeme})
	store.Insert(DomainEvent{Type: "new_user", Data: foobar})

	scanned := []Record{}

	err := store.Scan("build", func(r Record) error {
		scanned = append(scanned, r)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, scanned, 2)
	de0 := scanned[0].DomainEvent
	require.Equal(t, graeme, de0.Data)

	de1 := scanned[1].DomainEvent
	require.Equal(t, foobar, de1.Data)

	err = store.Scan("build", func(r Record) error {
		return errors.New("Should never get here")
	})
	require.NoError(t, err)

	gg := []byte{3}
	store.Insert(DomainEvent{Type: "new_user", Data: gg})

	var newScanned Record
	err = store.Scan("build", func(r Record) error {
		newScanned = r
		return nil
	})
	require.NoError(t, err)

	de2 := newScanned.DomainEvent
	require.Equal(t, gg, de2.Data)
}
