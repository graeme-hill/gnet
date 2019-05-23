package views

import (
	"context"
	"time"

	"github.com/graeme-hill/gnet/sys/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type DomainEvent struct {
	ID   int64
	Data []byte
	Date time.Time
	Type string
}

type ScanClient struct {
	client pb.DomainEventsClient
	conn   *grpc.ClientConn
}

func NewScanClient(addr string) (*ScanClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return &ScanClient{}, errors.Wrapf(err, "Cannot connect grpc client using addr '%s'", addr)
	}

	return &ScanClient{client: pb.NewDomainEventsClient(conn), conn: conn}, nil
}

type ScanHandler func(de DomainEvent) (bool, error)

func (sc *ScanClient) Scan(pointer uint32, handler ScanHandler) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Open two-way stream.
	stream, err := sc.client.Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to initiate scan stream")
	}

	// Send first message saying where to start scanning from.
	err = stream.Send(&pb.ScanRequest{
		Command: &pb.ScanRequest_ResumeCommand{
			ResumeCommand: &pb.ScanRequestResume{
				Pointer: pointer,
			},
		},
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

		onEvent := func(e *pb.ScanResponseDomainEvent) (bool, error) {
			// Delegate actual work to handler
			keepGoing, err := handler(DomainEvent{
				ID:   e.Id,
				Data: e.Data,
				Date: time.Unix(e.Date, 0),
				Type: e.Type,
			})
			if err != nil {
				return false, errors.Wrap(err, "Abandoning scan because a handler failed")
			}

			// Tell the server we good for this event.
			err = stream.Send(&pb.ScanRequest{
				Command: &pb.ScanRequest_StatusCommand{
					StatusCommand: &pb.ScanRequestStatus{
						LastReceived: e.Id,
					},
				},
			})
			if err != nil {
				return false, errors.Wrap(err, "Abandoning scan because send failed")
			}

			return keepGoing, nil
		}

		switch cmd := sr.Command.(type) {
		case *pb.ScanResponse_Event:
			keepGoing, err := onEvent(cmd.Event)
			if err != nil {
				return err
			}
			if !keepGoing {
				return nil
			}
		case *pb.ScanResponse_Complete:
			break
		default:
			return errors.New("unknown scan response command type")
		}
	}
}

func (sc *ScanClient) Close() {
	sc.conn.Close()
}
