package eventstore

import "time"

type ScanHandler = func(Record) error

type EventStore interface {
	Insert(DomainEvent) error
	Scan(pointer uint32, handler ScanHandler) error
}

type DomainEvent struct {
	Type string
	Date time.Time
	Data []byte
}

type Record struct {
	ID          int64
	DomainEvent DomainEvent
}

type Pointer struct {
	ID int64
}
