package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/stretchr/testify/require"

	"github.com/graeme-hill/gnet/sys/rpc-domainevents/pb"
	"google.golang.org/grpc"
)

func runServer(t *testing.T) {
	listen, err := net.Listen("tcp", ":50505")
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &server{
		store: eventstore.NewEventStoreConn("mem"),
	})

	if err := s.Serve(listen); err != nil {
		t.Errorf("failed to server: %v", err)
	}
}

func waitForServer(addr string, delay time.Duration, attempts int) (*grpc.ClientConn, error) {
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

func TestServer(t *testing.T) {
	go runServer(t)
	conn, err := waitForServer("localhost:50505", 100*time.Millisecond, 10)
	if err != nil {
		t.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewDomainEventsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data0 := []byte{1, 2}
	data1 := []byte{3, 4}

	_, err = c.InsertDomainEvent(ctx, &pb.InsertDomainEventRequest{
		Type: "foo",
		Data: data0,
	})
	require.NoError(t, err)

	_, err = c.InsertDomainEvent(ctx, &pb.InsertDomainEventRequest{
		Type: "foo",
		Data: data1,
	})
	require.NoError(t, err)

	stream, err := c.Scan(ctx)
	require.NoError(t, err)

	err = stream.Send(&pb.ScanRequest{
		Pointer: 6,
		After:   -1,
	})
	require.NoError(t, err)

	sr, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(0), sr.Id)
	require.Equal(t, data0, sr.Data)

	sr, err = stream.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), sr.Id)
	require.Equal(t, data1, sr.Data)
}
