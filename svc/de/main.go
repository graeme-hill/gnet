package main

import (
	"context"
	"log"
	"net"

	pb "github.com/graeme-hill/gnet/svc/de/proto"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) InsertDomainEvent(ctx context.Context, in *pb.InsertDomainEventRequest) (*pb.InsertDomainEventResponse, error) {
	return nil, &pb.InsertDomainEventResponse{}
}

func main() {
	listen, err := net.Listen("tcp", 50505)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
