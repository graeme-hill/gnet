package scanclient

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/rpc-domainevents/server"
	pb "github.com/graeme-hill/gnet/sys/scanclient/pbscanclient"
)

func waitForScanClient(addr string, delay time.Duration, attempts int) (*ScanClient, error) {
	var err error = nil
	for i := 0; i < attempts; i++ {
		time.Sleep(delay)
		client, err := NewScanClient(addr)
		if err == nil {
			return client, nil
		}
	}
	return nil, err
}

func TestScanClient(t *testing.T) {
	go func() {
		_ = server.RunServer(":50505", eventstore.NewEventStoreConn("mem"))
	}()

	scanClient, err := waitForScanClient("localhost:50505", 100*time.Millisecond, 10)
	require.NoError(t, err)
	require.NotNil(t, scanClient)
	defer scanClient.Close()

	c := pb.NewDomainEventsClient(scanClient.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	events := []DomainEvent{}
	err = scanClient.Scan(1, -1, func(de DomainEvent) error {
		events = append(events, de)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, events, 2)

}
