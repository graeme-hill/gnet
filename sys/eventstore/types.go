package eventstore

import "time"

type ScanHandler = func(Record) error

type EventStore interface {
	Insert(DomainEvent) error
	Scan(pointer uint32, handler ScanHandler) error
	GetPointer(pointer uint32) (int64, error)
	SetPointer(pointer uint32, lastHandled int64) error
}

type DomainEvent struct {
	Type string
	Date time.Time
	Data []byte
}

type Record struct {
	ID          int64
	DomainEvent DomainEvent
	Date        time.Time
}

type Pointer struct {
	ID int64
}
