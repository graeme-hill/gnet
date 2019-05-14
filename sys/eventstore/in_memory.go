package eventstore

import (
	"sync"

	"github.com/pkg/errors"
)

// InMemEventStore implements EventStore interface and just uses a slice.
type InMemEventStore struct {
	records  []Record
	mutex    *sync.Mutex
	pointers map[string]*pointer
}

type pointer struct {
	id    int64
	mutex sync.Mutex
}

func (e *InMemEventStore) requirePointer(key string) *pointer {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	p, found := e.pointers[key]
	if !found {
		p = &pointer{id: -1}
		e.pointers[key] = p
	}

	return p
}

func (e *InMemEventStore) setPointer(key string, id int64) {
	e.pointers[key].id = id
}

func (e *InMemEventStore) nextID() int64 {
	return int64(len(e.records))
}

// Insert adds a new event to the store.
func (e *InMemEventStore) Insert(de DomainEvent) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.records = append(e.records, Record{
		ID:          e.nextID(),
		DomainEvent: de,
	})

	return nil
}

// Scan iterates through events for the given pointer.
func (e InMemEventStore) Scan(scanKey string, handler ScanHandler) error {
	p := e.requirePointer(scanKey)
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, record := range e.records[p.id+1:] {
		err := handler(record)
		if err != nil {
			return errors.Wrap(err, "scan handler returned an error")
		}
		e.setPointer(scanKey, record.ID)
	}
	return nil
}

// NewInMemEventStore creates a totally new in-memory store.
func NewInMemEventStore() EventStore {
	return &InMemEventStore{
		pointers: map[string]*pointer{},
		mutex:    &sync.Mutex{},
	}
}
