package scan

import (
	"context"
	"time"

	"github.com/graeme-hill/gnet/sys/scan/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type DomainEvent struct {
	ID   int64
	Data []byte
	Date time.Time
}

type ScanClient struct {
	client pb.DomainEventsClient
}

func NewScanClient(addr string) (*ScanClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return &ScanClient{}, errors.Wrapf(err, "Cannot connect grpc client using addr '%s'", addr)
	}

	return &ScanClient{client: pb.NewDomainEventsClient(conn)}, nil
}

type ScanHandler func(de DomainEvent) error

func (sc *ScanClient) Scan(pointer uint32, after int64, handler ScanHandler) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Open two-way stream.
	stream, err := sc.client.Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to initiate scan stream")
	}

	// Send first message saying where to start scanning from.
	err = stream.Send(&pb.ScanRequest{
		Pointer: pointer,
		After:   after,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to send scan request")
	}

	for {
		// Get next event.
		sr, err := stream.Recv()
		if err != nil {
			return errors.Wrap(err, "Failed to recv next event")
		}

		// Delegate actual work to handler
		err = handler(DomainEvent{
			ID:   sr.Id,
			Data: sr.Data,
			Date: time.Unix(sr.Date, 0),
		})
		if err != nil {
			return errors.Wrap(err, "Abandoning scan because a handler failed")
		}

		// Tell the server we good for this event.
		stream.Send(&pb.ScanRequest{
			Pointer: pointer,
			After:   sr.Id,
		})
	}
}
