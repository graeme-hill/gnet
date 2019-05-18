package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/rpc-domainevents/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type server struct {
	store eventstore.EventStore
}

func (s *server) InsertDomainEvent(ctx context.Context, in *pb.InsertDomainEventRequest) (*pb.InsertDomainEventResponse, error) {
	err := s.store.Insert(eventstore.DomainEvent{
		Type: in.Type,
		Date: time.Now(),
		Data: in.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.InsertDomainEventResponse{}, nil
}

func (s *server) Scan(stream pb.DomainEvents_ScanServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Failed to receive initial message from client")
	}

	err = s.store.Scan(req.Pointer, func(rec eventstore.Record) error {
		err = stream.Send(&pb.ScanResponse{
			Id:   rec.ID,
			Data: rec.DomainEvent.Data,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to send domain event to client")
		}
		return nil
	})

	return err
}

func main() {
	listen, err := net.Listen("tcp", "50505")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &server{
		store: eventstore.NewEventStoreConn("mem"),
	})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
