package eventstore

import (
	"sync"

	"github.com/pkg/errors"
)

var stores = map[string]*InMemEventStore{}

type InMemEventStore struct {
	records  []Record
	mutex    *sync.Mutex
	pointers map[uint32]*pointer
}

type pointer struct {
	id    int64
	mutex sync.Mutex
}

func (e *InMemEventStore) GetPointer(pointer uint32) (int64, error) {
	p, ok := e.pointers[pointer]
	if !ok {
		return 0, errors.Errorf("pointer %d does not exist", pointer)
	}
	return p.id, nil
}

func (e *InMemEventStore) SetPointer(pointer uint32, lastHandled int64) error {
	p, ok := e.pointers[pointer]
	if !ok {
		return errors.Errorf("pointer %d does not exist", pointer)
	}
	p.id = lastHandled
	return nil
}

func (e *InMemEventStore) requirePointer(pointerID uint32) *pointer {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	p, found := e.pointers[pointerID]
	if !found {
		p = &pointer{id: -1}
		e.pointers[pointerID] = p
	}

	return p
}

func (e *InMemEventStore) setPointer(pointerID uint32, id int64) {
	e.pointers[pointerID].id = id
}

func (e *InMemEventStore) nextID() int64 {
	return int64(len(e.records))
}

func (e *InMemEventStore) Insert(de DomainEvent) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.records = append(e.records, Record{
		ID:          e.nextID(),
		DomainEvent: de,
	})

	return nil
}

func (e InMemEventStore) Scan(pointer uint32, handler ScanHandler) error {
	p := e.requirePointer(pointer)
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, record := range e.records[p.id+1:] {
		err := handler(record)
		if err != nil {
			return errors.Wrap(err, "scan handler returned an error")
		}
		e.setPointer(pointer, record.ID)
	}
	return nil
}

func NewInMemEventStore(connStr string) EventStore {
	store, ok := stores[connStr]
	if !ok {
		store = &InMemEventStore{
			pointers: map[uint32]*pointer{},
			mutex:    &sync.Mutex{},
		}
		stores[connStr] = store
	}
	return store
}
