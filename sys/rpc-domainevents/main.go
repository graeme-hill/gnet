package main

import (
	"log"

	"github.com/graeme-hill/gnet/sys/eventstore"
	"github.com/graeme-hill/gnet/sys/rpc-domainevents/server"
)

func main() {
	err := server.RunServer(":50505", eventstore.NewEventStoreConn("mem"))
	if err != nil {
		log.Printf("Server closed with error: %v\n", err)
	}
}
