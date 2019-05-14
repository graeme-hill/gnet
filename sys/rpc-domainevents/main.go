package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/graeme-hill/gnet/lib/eventstore"
	pb "github.com/graeme-hill/gnet/svc/de/proto"
	"google.golang.org/grpc"
	"github.com/pkg/errors"
)

type server struct {
	store *eventstore.InMemEventStore
}

func (s *server) InsertDomainEvent(ctx context.Context, in *pb.InsertDomainEventRequest) (*pb.InsertDomainEventResponse, error) {
	return &pb.InsertDomainEventResponse{}, nil
}

func (s *server) Scan(stream pb.DomainEvents_ScanServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Failed to receive initial message from client")
	}

	s.store.Scan("TODO", func(rec eventstore.Record) error {
		err = stream.Send(&pb.ScanResponse{
			Id: rec.Id,
			Data: rec.DomainEvent.Data,
		})
		if err != nil {
			return errors.Wrap(err, "Failed to send domain event to client")
		}
	})

	return nil
}

func main() {
	listen, err := net.Listen("tcp", "50505")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &server{
		store: eventstore.NewInMemEventStore(),
	})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
