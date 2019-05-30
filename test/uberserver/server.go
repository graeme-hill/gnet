package uberserver

import (
	"context"
	"log"

	photoserver "github.com/graeme-hill/gnet/photos/web-photos/server"
	"github.com/graeme-hill/gnet/sys/gnet"
	deserver "github.com/graeme-hill/gnet/sys/rpc-domainevents/server"
)

type UberServer struct {
	Connections  gnet.Connections
	deErrChan    <-chan error
	photoErrChan <-chan error
	over         chan map[string]error
}

func (u *UberServer) Done() <-chan map[string]error {
	if u.over == nil {
		u.over = make(chan map[string]error)

		go func() {
			log.Println("AGG: waiting")
			u.over <- map[string]error{
				"rpc-domainevents": <-u.deErrChan,
				"web-photos":       <-u.photoErrChan,
			}
			log.Println("AGG: doned")
			close(u.over)
		}()
	}

	return u.over
}

func StartUberServer(ctx context.Context) UberServer {
	connections := gnet.Connections{
		EventStore:      ":memory:",
		FileStore:       ":memory:",
		KeyValueStore:   ":memory:",
		PhotosWebAPI:    "http://localhost:8000",
		DomainEventsRPC: "localhost:9000",
	}

	return UberServer{
		Connections: connections,
		deErrChan: deserver.Run(ctx, deserver.Options{
			Addr:              connections.DomainEventsRPC,
			EventStoreConnStr: connections.EventStore,
		}),
		photoErrChan: photoserver.Run(ctx, photoserver.Options{
			Addr:                connections.PhotosWebAPI,
			DomainEventsRPCAddr: connections.DomainEventsRPC,
			EventStoreConnStr:   connections.EventStore,
			FileStoreConnStr:    connections.FileStore,
		}),
	}
}
