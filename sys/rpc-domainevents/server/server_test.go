package server

import (
	"context"
	"testing"
	"time"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/stretchr/testify/require"

	"github.com/graeme-hill/gnet/sys/pb"
)

func TestServer(t *testing.T) {
	go func() {
		_ = RunServer(":50505", eventstore.NewEventStoreConn("mem"))
	}()
	conn, err := WaitForServer("localhost:50505", 100*time.Millisecond, 10)
	require.NoError(t, err)
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
		Command: &pb.ScanRequest_ResumeCommand{
			ResumeCommand: &pb.ScanRequestResume{
				Pointer: 6,
			},
		},
	})
	require.NoError(t, err)

	sr, err := stream.Recv()
	require.NoError(t, err)
	eventCmd, ok := sr.Command.(*pb.ScanResponse_Event)
	require.True(t, ok)
	require.Equal(t, int64(0), eventCmd.Event.Id)
	require.Equal(t, data0, eventCmd.Event.Data)

	sr, err = stream.Recv()
	require.NoError(t, err)
	eventCmd, ok = sr.Command.(*pb.ScanResponse_Event)
	require.Equal(t, int64(1), eventCmd.Event.Id)
	require.Equal(t, data1, eventCmd.Event.Data)

	sr, err = stream.Recv()
	require.NoError(t, err)
	_, ok = sr.Command.(*pb.ScanResponse_Complete)
	require.True(t, ok)
}
