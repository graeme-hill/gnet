package eventstore

import (
	"testing"
	"time"

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
	store := NewInMemEventStore("mem")
	graeme := []byte{1}
	foobar := []byte{2}
	now := time.Now()
	err := store.Insert(DomainEvent{Type: "new_user", Data: graeme, Date: now})
	require.NoError(t, err)
	err = store.Insert(DomainEvent{Type: "new_user", Data: foobar, Date: now})
	require.NoError(t, err)

	scanned := []Record{}

	err = store.Scan(1, func(r Record) error {
		scanned = append(scanned, r)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, scanned, 2)
	de0 := scanned[0].DomainEvent
	require.Equal(t, graeme, de0.Data)
	require.Equal(t, now, de0.Date)

	de1 := scanned[1].DomainEvent
	require.Equal(t, foobar, de1.Data)
	require.Equal(t, now, de1.Date)

	err = store.Scan(1, func(r Record) error {
		return errors.New("Should never get here")
	})
	require.NoError(t, err)

	gg := []byte{3}
	err = store.Insert(DomainEvent{Type: "new_user", Data: gg, Date: now})
	require.NoError(t, err)

	var newScanned Record
	err = store.Scan(1, func(r Record) error {
		newScanned = r
		return nil
	})
	require.NoError(t, err)

	de2 := newScanned.DomainEvent
	require.Equal(t, gg, de2.Data)
	require.Equal(t, now, de2.Date)
}
