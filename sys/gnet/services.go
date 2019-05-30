package gnet

import "sync"

type Connections struct {
	EventStore      string
	FileStore       string
	KeyValueStore   string
	PhotosWebAPI    string
	DomainEventsRPC string
}

type Service struct {
	overChan    <-chan error
	runningChan <-chan struct{}
	done        bool
	err         error
	running     bool
	mu          sync.Mutex
}

func (s *Service) Started() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return s.runningChan
	}

	alreadyStartedChan := make(chan struct{}, 1)
	alreadyStartedChan <- struct{}{}
	return alreadyStartedChan
}

func (s *Service) Done() <-chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.done {
		return s.overChan
	}

	alreadyDoneChan := make(chan error, 1)
	alreadyDoneChan <- s.err
	return alreadyDoneChan
}

func NewService(over <-chan error, running <-chan struct{}) *Service {
	return &Service{
		overChan:    over,
		runningChan: running,
		done:        false,
		running:     false,
		err:         nil,
	}
}
