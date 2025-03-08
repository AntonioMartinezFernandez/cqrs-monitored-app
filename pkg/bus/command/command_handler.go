package command

import (
	"context"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type CommandHandler interface {
	Handle(ctx context.Context, command bus.Dto) error
}
