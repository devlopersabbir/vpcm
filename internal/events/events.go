package events

import (
	"sync"
)

type Event struct {
	Type    string
	Payload any
}

type HandlerFunc func(Event)

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]HandlerFunc
}

var globalBus = NewEventBus()

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]HandlerFunc),
	}
}

func Publish(event Event) {
	globalBus.Publish(event)
}

func Subscribe(eventType string, handler HandlerFunc) {
	globalBus.Subscribe(eventType, handler)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, exists := eb.subscribers[event.Type]
	if !exists {
		return
	}

	for _, handler := range handlers {
		go handler(event)
	}
}

func (eb *EventBus) Subscribe(eventType string, handler HandlerFunc) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}
