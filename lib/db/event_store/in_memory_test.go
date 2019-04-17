package event_store

import (
	"testing"

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
	store.Insert(NewUserEvent{name: "graeme"})
	store.Insert(NewUserEvent{name: "foobar"})

	scanned := []Record{}

	err := store.ScanSince(-1, func(r Record) error {
		scanned = append(scanned, r)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, scanned, 2)
	de0, isNewUserEvent := scanned[0].DomainEvent.(NewUserEvent)
	require.True(t, isNewUserEvent)
	require.Equal(t, "graeme", de0.name)

	de1, isNewUserEvent := scanned[1].DomainEvent.(NewUserEvent)
	require.True(t, isNewUserEvent)
	require.Equal(t, "foobar", de1.name)
}
