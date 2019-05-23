package builders

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/graeme-hill/gnet/photos/events"
	"github.com/graeme-hill/gnet/photos/pb"
	"github.com/graeme-hill/gnet/sys/keyvaluestore"
	"github.com/graeme-hill/gnet/sys/views"
)

func buildPhotoMetadata(de views.DomainEvent) error {
	if de.Type == events.EventPhotoUploaded {
		event := pb.PhotoUploaded{}
		err := proto.Unmarshal(de.Data, &event)

		if err != nil {
			log.Println("failed to deserialize domain event in buildPhotoMetadata")
			return nil
		}

		kvs := keyvaluestore.NewKeyValueStoreConn(":memory:")
		err = kvs.Set(fmt.Sprintf("ph:%s", event.Path), []byte{})
		if err != nil {
			log.Println("failed to write photo metadata to database in buildPhotoMetadata")
			return nil
		}
	}

	return nil
}

type PhotoMetadataBuilder struct{}

func (b *PhotoMetadataBuilder) Key() uint32 {
	return 6
}

func (b *PhotoMetadataBuilder) Types() []string {
	return []string{events.EventPhotoUploaded}
}

func (b *PhotoMetadataBuilder) OnDomainEvent(de views.DomainEvent) error {
	return buildPhotoMetadata(de)
}

func newPhotoMetadataViewWorker() *views.Worker {
	return views.NewWorker(&PhotoMetadataBuilder{}, "localhost:50505")
}
