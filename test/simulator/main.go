package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/graeme-hill/gnet/test/fakeuploader"
	"github.com/graeme-hill/gnet/test/uberserver"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	uberServer := uberserver.StartUberServer(ctx)

	fuOver := fakeuploader.Run(ctx, uberServer.Connections.PhotosWebAPI)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	cancelRequested := false
loop:
	for {
		select {
		case errs := <-uberServer.Done():
			// stop simulation jobs
			fuErr := <-fuOver
			if fuErr == nil {
				log.Print("Stopped fakeuploader: OK :)")
			} else {
				log.Printf("Stopped fakeuploader: %v", fuErr)
			}

			// stop services
			for service, err := range errs {
				if err != nil {
					log.Printf("Error shutting down %s: %v", service, err)
				}
			}
			break loop

		case <-sigChan: // CTRL-C
			log.Println("Stopping services...")
			cancel()
			cancelRequested = true

		case <-time.After(10 * time.Second):
			if cancelRequested {
				log.Println("TIMEOUT")
				break loop
			}
		}
	}
}

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println("finished one")
// 				return
// 			default:
// 			}
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println("finished two")
// 				return
// 			default:
// 			}
// 		}
// 	}()

// 	fmt.Println("waiting 10 seconds")
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("cancel")
// 	cancel()
// 	fmt.Println("waiting 2 seconds")
// 	time.Sleep(2 * time.Second)
// 	fmt.Println("done")
// }
