package event_store

import (
	"github.com/pkg/errors"
	"sync"
)

type inMemEventStore struct {
	records []Record
	mutex   sync.Mutex
}

func (e *inMemEventStore) nextID() int {
	return len(e.records)
}

func (e *inMemEventStore) Insert(de DomainEvent) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.records = append(e.records, Record{
		ID:          e.nextID(),
		DomainEvent: de,
	})

	return nil
}

func (e inMemEventStore) ScanSince(id int, handler ScanHandler) error {
	for _, record := range e.records[id+1:] {
		err := handler(record)
		if err != nil {
			return errors.Wrap(err, "scan handler returned an error")
		}
	}
	return nil
}

func NewInMemEventStore() inMemEventStore {
	return inMemEventStore{}
}
