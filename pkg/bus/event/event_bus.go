package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

// EventBus is the central component that handles event dispatching and subscription.
type EventBus struct {
	mut         *sync.RWMutex
	subscribers map[string][]EventHandler
}

// NewEventBus creates a new instance of EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		mut:         &sync.RWMutex{},
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe adds a new handler for a specific event.
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	eb.mut.Lock()
	defer eb.mut.Unlock()

	// Append the handler to the list of handlers for the given event name
	eb.subscribers[eventName] = append(eb.subscribers[eventName], handler)
}

// Unsubscribe removes a handler for a specific event.
func (eb *EventBus) Unsubscribe(eventName string, handler EventHandler) {
	eb.mut.Lock()
	defer eb.mut.Unlock()

	handlers, ok := eb.subscribers[eventName]
	if !ok {
		return
	}

	// Remove the handler
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			eb.subscribers[eventName] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// PublishOne publishes an event to all subscribers of that event.
func (eb *EventBus) PublishOne(ctx context.Context, event bus.Event) {
	eb.mut.RLock()
	defer eb.mut.RUnlock()

	// Check if there are any subscribers for this event
	handlers, ok := eb.subscribers[event.Name()]
	if !ok {
		return
	}

	// Call each handler with the event
	for _, handler := range handlers {
		go handler.Handle(ctx, event)
	}
}

func (eb *EventBus) Publish(ctx context.Context, events bus.Events) {
	eb.mut.RLock()
	defer eb.mut.RUnlock()

	for _, event := range events {
		eb.PublishOne(ctx, event)
	}
}
