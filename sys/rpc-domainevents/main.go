package main

import (
	"context"
	"log"

	"github.com/graeme-hill/gnet/sys/rpc-domainevents/server"
)

func main() {
	ctx := context.Background()
	errChan := server.Run(ctx, server.Options{
		Addr:              ":50505",
		EventStoreConnStr: ":memory:",
	})

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Server closed with error: %v\n", err)
		} else {
			log.Printf("Server closed cleanly")
		}
	}
}
