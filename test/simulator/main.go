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
				log.Print("OFFLINE fakeuploader: OK :)")
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
