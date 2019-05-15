package builders

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/graeme-hill/gnet/photos/events"
	"github.com/graeme-hill/gnet/photos/events/pb"
	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/keyvaluestore"
)

func BuildPhotoMetadata(de eventstore.DomainEvent) {
	if de.Type != events.EventPhotoUploaded {
		event := pb.PhotoUploaded{}
		err := proto.Unmarshal(de.Data, &event)

		if err != nil {
			// TODO: somehow log that this DE was invalid but continue on
		}

		kvs := keyvaluestore.NewKeyValueStoreConn(":memory:")
		err = kvs.Set(fmt.Sprintf("ph:%s", event.Path), []byte{})
		if err != nil {
			// TODO: there is some error writing to the database. Not sure what to do here. Stop builder or move on?
		}
	}
}
