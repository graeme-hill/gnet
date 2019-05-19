package server

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

type Server struct {
	store eventstore.EventStore
}

func (s *Server) InsertDomainEvent(ctx context.Context, in *pb.InsertDomainEventRequest) (*pb.InsertDomainEventResponse, error) {
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

func (s *Server) Scan(stream pb.DomainEvents_ScanServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "Failed to receive initial message from client")
	}

	go func() {
		for {
			_, err2 := stream.Recv()
			if err2 != nil {
				log.Printf("server failed to recv '%v'\n", err2)
			} else {
				log.Println("server succesfully recv'd")
			}
		}
	}()

	go func() {
		_ = s.store.Scan(req.Pointer, func(rec eventstore.Record) error {
			err3 := stream.Send(&pb.ScanResponse{
				Id:   rec.ID,
				Data: rec.DomainEvent.Data,
			})
			if err3 != nil {
				return errors.Wrap(err, "Failed to send domain event to client")
			}
			return nil
		})
	}()

	return nil
}

func RunServer(addr string, estore eventstore.EventStore) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &Server{
		store: estore,
	})

	return s.Serve(listen)
}

func WaitForServer(addr string, delay time.Duration, attempts int) (*grpc.ClientConn, error) {
	var err error = nil
	for i := 0; i < attempts; i++ {
		time.Sleep(delay)
		conn, err := grpc.Dial("localhost:50505", grpc.WithInsecure())
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}
