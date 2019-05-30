package eventstore

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func NewEventStoreConn(connStr string) EventStore {
	// when there is a non-memory implementation then make this smarter
	return NewInMemEventStore(connStr)
}

func NewDomainEvent(eventType string, msg proto.Message) (DomainEvent, error) {
	payload, err := proto.Marshal(msg)

	if err != nil {
		return DomainEvent{}, errors.Wrapf(err, "Failed to marshal %s domain event", eventType)
	}

	return DomainEvent{
		Type: eventType,
		Data: payload,
		Date: time.Now(),
	}, nil
}
