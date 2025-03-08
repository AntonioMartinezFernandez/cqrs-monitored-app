package query

import (
	"context"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type QueryHandler interface {
	Handle(ctx context.Context, query bus.Dto) (any, error)
}
