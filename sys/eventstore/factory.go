package eventstore

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var NewEventStoreConn func(string) EventStore = NewInMemEventStore

func NewDomainEvent(eventType string, msg proto.Message) (DomainEvent, error) {
	payload, err := proto.Marshal(msg)

	if err != nil {
		return DomainEvent{}, errors.Wrapf(err, "Failed to marshal %s domain event", eventType)
	}

	return DomainEvent{
		Type: eventType,
		Data: payload,
	}, nil
}
