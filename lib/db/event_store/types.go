package event_store

type ScanHandler = func(Record) error

type EventStore interface {
	Insert(DomainEvent) error
	ScanSince(id int, handler ScanHandler) error
}

type DomainEvent interface {
	Type() string
}

type Record struct {
	ID          int
	DomainEvent DomainEvent
}
