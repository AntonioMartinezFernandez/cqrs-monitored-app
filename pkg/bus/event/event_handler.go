package event

import (
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type EventHandler interface {
	Handle(event bus.Event) error
}
