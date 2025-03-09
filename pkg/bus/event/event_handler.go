package event

import (
	"context"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type EventHandler interface {
	Handle(ctx context.Context, event bus.Event)
}
