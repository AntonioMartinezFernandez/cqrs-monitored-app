package domain

import (
	"sync"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type EventRecorder struct {
	mut    sync.Mutex
	events bus.Events
}

func NewEventRecorder() *EventRecorder {
	return &EventRecorder{
		mut:    sync.Mutex{},
		events: bus.Events{},
	}
}

func (er *EventRecorder) Record(event bus.Event) {
	er.mut.Lock()
	defer er.mut.Unlock()

	er.events = append(er.events, event)
}

func (er *EventRecorder) Flush() {
	er.mut.Lock()
	defer er.mut.Unlock()

	er.events = make(bus.Events, 0)
}

func (er *EventRecorder) Pull() bus.Events {
	events := er.events

	er.Flush()

	return events
}
