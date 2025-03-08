package query

import (
	"context"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"

	"sync"
)

type Bus interface {
	RegisterQuery(query bus.Dto, handler QueryHandler) error
	Ask(ctx context.Context, dto bus.Dto) (any, error)
}

type QueryBus struct {
	handlers map[string]QueryHandler
	lock     sync.Mutex
	logger   pkg_logger.Logger
}

func InitQueryBus(logger pkg_logger.Logger) *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler, 0),
		lock:     sync.Mutex{},
		logger:   logger,
	}
}

type QueryAlreadyRegistered struct {
	message   string
	queryName string
}

func (i QueryAlreadyRegistered) Error() string {
	return i.message
}

func NewQueryAlreadyRegistered(message string, queryName string) QueryAlreadyRegistered {
	return QueryAlreadyRegistered{message: message, queryName: queryName}
}

type QueryNotRegistered struct {
	message string
}

func (i QueryNotRegistered) Error() string {
	return i.message
}

func NewQueryNotRegistered(message string, queryName string) QueryAlreadyRegistered {
	return QueryAlreadyRegistered{message: message, queryName: queryName}
}

func (bus *QueryBus) RegisterQuery(query bus.Dto, handler QueryHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	queryName := query.Type()

	if _, ok := bus.handlers[queryName]; ok {
		return NewQueryAlreadyRegistered("query already registered", queryName)
	}

	bus.handlers[queryName] = handler

	return nil
}

func (bus *QueryBus) Ask(ctx context.Context, query bus.Dto) (any, error) {
	queryName := query.Type()

	if handler, ok := bus.handlers[queryName]; ok {
		response, err := bus.doAsk(ctx, handler, query)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, NewQueryNotRegistered("query not registered", queryName)
}

func (bus *QueryBus) doAsk(ctx context.Context, handler QueryHandler, query bus.Dto) (any, error) {
	return handler.Handle(ctx, query)
}

type QueryNotValid struct {
	message string
}

func (i QueryNotValid) Error() string {
	return i.message
}
