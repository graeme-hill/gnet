package events

import (
	"github.com/graeme-hill/gnet/sys/eventstore"
)

func NewPhotoUploadedEvent(path string) eventstore.DomainEvent {
	return eventstore.DomainEvent{
		Type: "photo_uploaded",
	}
}
