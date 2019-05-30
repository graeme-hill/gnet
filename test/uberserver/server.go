package uberserver

import (
	"context"
	"log"

	photoserver "github.com/graeme-hill/gnet/photos/web-photos/server"
	"github.com/graeme-hill/gnet/sys/gnet"
	deserver "github.com/graeme-hill/gnet/sys/rpc-domainevents/server"
)

type UberServer struct {
	Connections           gnet.Connections
	domainEventRPCService gnet.Service
	photosAPIService      gnet.Service
	over                  chan map[string]error
	running               chan struct{}
}

func (u *UberServer) wait() {
	_ = <-u.domainEventRPCService.Running
	_ = <-u.photosAPIService.Running
}

func (u *UberServer) Done() <-chan map[string]error {
	if u.over == nil {
		u.over = make(chan map[string]error)

		go func() {
			log.Println("AGG: waiting")
			u.over <- map[string]error{
				"rpc-domainevents": <-u.domainEventRPCService.Over,
				"web-photos":       <-u.photosAPIService.Over,
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

	result := UberServer{
		Connections: connections,
		domainEventRPCService: deserver.Run(ctx, deserver.Options{
			Addr:              connections.DomainEventsRPC,
			EventStoreConnStr: connections.EventStore,
		}),
		photosAPIService: photoserver.Run(ctx, photoserver.Options{
			Addr:                connections.PhotosWebAPI,
			DomainEventsRPCAddr: connections.DomainEventsRPC,
			EventStoreConnStr:   connections.EventStore,
			FileStoreConnStr:    connections.FileStore,
		}),
	}

	result.wait()
	return result
}
