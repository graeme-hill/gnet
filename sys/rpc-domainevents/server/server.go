package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	eventStoreConnStr string
}

func (s *Server) InsertDomainEvent(ctx context.Context, in *pb.InsertDomainEventRequest) (*pb.InsertDomainEventResponse, error) {
	store := eventstore.NewEventStoreConn(s.eventStoreConnStr)

	err := store.Insert(eventstore.DomainEvent{
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

	resumeCommand, ok := req.Command.(*pb.ScanRequest_ResumeCommand)
	if !ok {
		return errors.New("First message from client must be resume command")
	}
	pointer := resumeCommand.ResumeCommand.Pointer

	finalID := int64(-1)
	receiverDone := make(chan struct{})
	senderDone := make(chan struct{})

	go func() {
		for {
			select {
			case <-stream.Context().Done():
				receiverDone <- struct{}{}
				return
			default:
			}
			msg, err2 := stream.Recv()
			if err2 != nil {
				log.Printf("server failed to recv '%v'\n", err2)
				break
			} else {
				log.Println("server succesfully recv'd")
				switch cmd := msg.Command.(type) {
				case *pb.ScanRequest_ResumeCommand:

				case *pb.ScanRequest_StatusCommand:
					if finalID >= 0 && cmd.StatusCommand.LastReceived >= finalID {
						log.Printf("complete")
						receiverDone <- struct{}{}
						return
					}
				}
			}
		}
	}()

	go func() {
		maxID := int64(-1)
		store := eventstore.NewEventStoreConn(s.eventStoreConnStr)
		_ = store.Scan(pointer, func(rec eventstore.Record) error {
			err3 := stream.Send(&pb.ScanResponse{
				Command: &pb.ScanResponse_Event{
					Event: &pb.ScanResponseDomainEvent{
						Id:   rec.ID,
						Data: rec.DomainEvent.Data,
						Date: rec.DomainEvent.Date.Unix(),
						Type: rec.DomainEvent.Type,
					},
				},
			})
			maxID = rec.ID
			if err3 != nil {
				return errors.Wrap(err, "Failed to send domain event to client")
			}
			return nil
		})

		completeErr := stream.Send(&pb.ScanResponse{
			Command: &pb.ScanResponse_Complete{
				Complete: &pb.ScanResponseComplete{},
			},
		})
		if completeErr != nil {
			log.Printf("failed to send complete message to client: %v", completeErr)
		}

		log.Printf("server finished scanning but waiting for client ack")
		finalID = maxID
		senderDone <- struct{}{}
	}()

	<-receiverDone
	<-senderDone

	return nil
}

func Run(ctx context.Context, opt Options) <-chan error {
	log.Println("RPC SERVER: listening " + opt.Addr)
	listen, err := net.Listen("tcp", opt.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDomainEventsServer(s, &Server{
		eventStoreConnStr: opt.EventStoreConnStr,
	})

	over := make(chan error)
	go start(ctx, s, listen, over)

	return over
}

func start(
	ctx context.Context,
	server *grpc.Server,
	listener net.Listener,
	over chan<- error,
) {
	go func() {
		over <- server.Serve(listener)
		log.Println("RPC SERVER: shutted down")
	}()

	log.Println("RPC SERVER: waiting for ctx to be done")
	select {
	case <-ctx.Done():
		log.Println("RPC SERVER: ctx done!")
		server.GracefulStop()
	}
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

type Options struct {
	EventStoreConnStr string
	Addr              string
}
