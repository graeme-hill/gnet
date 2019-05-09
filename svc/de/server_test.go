package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/graeme-hill/gnet/svc/de/proto"
	"google.golang.org/grpc"
)

func runServer(t *testing.T) {
	listen, err := net.Listen("tcp", ":50505")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		t.Fatalf("failed to server: %v", err)
	}
}

func TestServer2(t *testing.T) {
	go runServer(t)
	time.Sleep(100, * time.Millisecond)

	client := eventstore.ScanClient("localhost:50505")
	client.Scan(1, 0, func
}

func TestServer(t *testing.T) {
	go runServer(t)
	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.Dial("localhost:50505", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewDomainEventsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.InsertDomainEvent(ctx, &pb.InsertDomainEventRequest{
		Type: "foo",
		Data: []byte{},
	})

	if err != nil {
		t.Fatalf("error inserting domain event: %v", err)
	}

	stream, err := c.Scan(ctx)
	if err != nil {
		t.Fatalf("failed to open scan stream: %v", err)
	}

	err = stream.Send(&pb.ScanRequest{
		Pointer: 6,
		After:   77,
	})
	if err != nil {
		t.Fatalf("failed to send scan req: %v", err)
	}

	sr, err := stream.Recv()
	if err != nil {
		t.Fatalf("failed to recv: %v", err)
	}

	if sr.Id != 77 {
		log.Fatalf("bad sr: %v", sr)
	}
}
