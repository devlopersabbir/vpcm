package events

import (
	"sync"
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	bus := NewEventBus()
	var wg sync.WaitGroup
	wg.Add(1)

	var received Event
	bus.Subscribe("TestEvent", func(ev Event) {
		received = ev
		wg.Done()
	})

	bus.Publish(Event{
		Type:    "TestEvent",
		Payload: "Hello World",
	})

	// Wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event delivery")
	}

	if received.Payload.(string) != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", received.Payload)
	}
}
