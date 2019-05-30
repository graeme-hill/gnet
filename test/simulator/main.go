package main

import (
	"context"
	"log"
	"os"
	"os/signal"

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
				if err == nil {
					log.Printf("Stopped %s: OK :)", service)
				} else {
					log.Printf("Stopped %s: Error: %v", service, err)
				}
			}
			break loop

		case <-sigChan: // CTRL-C
			log.Println("Stopping services...")
			cancel()
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
