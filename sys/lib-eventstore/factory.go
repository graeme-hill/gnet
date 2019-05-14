package eventstore

var NewEventStoreConn func() EventStore = NewInMemEventStore
