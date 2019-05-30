package gnet

type Connections struct {
	EventStore      string
	FileStore       string
	KeyValueStore   string
	PhotosWebAPI    string
	DomainEventsRPC string
}

type Service struct {
	Over    <-chan error
	Running <-chan struct{}
}
