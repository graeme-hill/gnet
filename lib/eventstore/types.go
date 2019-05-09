package eventstore

type ScanHandler = func(Record) error

type EventStore interface {
	Insert(DomainEvent) error
	Scan(pointer string, handler ScanHandler) error
}

type DomainEvent interface {
	Type() string
}

type Record struct {
	ID          int
	DomainEvent DomainEvent
}

type Pointer struct {
	ID int
}
