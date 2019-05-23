package views

import (
	"log"
	"time"
)

type Worker struct {
	builder             Builder
	stop                chan struct{}
	stopped             chan struct{}
	receivedStopMessage bool
	scanClient          *ScanClient
	eventsAddr          string
}

func NewWorker(b Builder, addr string) *Worker {
	return &Worker{
		builder:             b,
		stop:                make(chan struct{}),
		stopped:             make(chan struct{}),
		receivedStopMessage: false,
		scanClient:          nil,
		eventsAddr:          addr,
	}
}

func (w *Worker) requireScanClient() *ScanClient {
	for !w.shouldStop() {
		scanClient, err := NewScanClient(w.eventsAddr)
		if err != nil {
			log.Printf("cannot connect to '%s': %s", w.eventsAddr, err.Error())
		} else {
			return scanClient
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (w *Worker) Start() {
	go w.loop()
}

// Tell the loop to stop looping and then wait for it to actually stop.
func (w *Worker) Stop() {
	w.stop <- struct{}{}
	<-w.stopped
}

// Keep processing domain events until told to stop.
func (w *Worker) loop() {
	// TODO: instead of loop and sleep to connection to the event server should just
	// stay open so it can push new events as they appear.
	for {
		w.doWork()
		if w.shouldStop() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	w.stopped <- struct{}{}
}

// Actually check for and handle a set of domain events.
func (w *Worker) doWork() {
	client := w.requireScanClient()
	if client == nil {
		// This probably means that it stopped trying to connect because the
		// worker was told to stop.
		return
	}

	err := client.Scan(w.builder.Key(), func(de DomainEvent) (bool, error) {
		deErr := w.builder.OnDomainEvent(DomainEvent)
		if deErr != nil {
			log.Printf("Error handling domain event: %v", deErr)
		}
		return !w.shouldStop(), nil
	})
}

// Check if the worker goroutine has been told to stop
func (w *Worker) shouldStop() bool {
	// If already received stop message then this worker is done and should
	// not do more work or expect more stop messages.
	if w.receivedStopMessage {
		return true
	}

	select {
	case <-w.stop:
		w.receivedStopMessage = true
		return true
	default:
		return false
	}
}
