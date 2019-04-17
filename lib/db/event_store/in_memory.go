package event_store

import (
	"github.com/pkg/errors"
	"sync"
)

type inMemEventStore struct {
	records  []Record
	mutex    sync.Mutex
	pointers map[string]*pointer
}

type pointer struct {
	id    int
	mutex sync.Mutex
}

func (e *inMemEventStore) requirePointer(key string) *pointer {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	p, found := e.pointers[key]
	if !found {
		p = &pointer{id: -1}
		e.pointers[key] = p
	}

	return p
}

func (e *inMemEventStore) setPointer(key string, id int) {
	e.pointers[key].id = id
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

func (e inMemEventStore) Scan(scanKey string, handler ScanHandler) error {
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

func NewInMemEventStore() inMemEventStore {
	return inMemEventStore{
		pointers: map[string]*pointer{},
	}
}
