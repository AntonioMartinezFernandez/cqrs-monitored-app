package command

import (
	"context"
	"log/slog"
	"reflect"
	"sync"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	mutex "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/mutex-service"
)

type Bus interface {
	RegisterCommand(command bus.Dto, handler CommandHandler) error
	Exec(ctx context.Context, dto bus.Dto) error
	ExecAsync(ctx context.Context, dto bus.Dto) error
	ReprocessFailedAsyncCommands(ctx context.Context, maxTimes int)
}

type CommandBus struct {
	handlers       map[string]CommandHandler
	mut            sync.Mutex
	logger         pkg_logger.Logger
	failedCommands chan *FailedCommand

	mutex mutex.MutexService
}

func InitCommandBus(logger pkg_logger.Logger, mutex mutex.MutexService) *CommandBus {
	return &CommandBus{
		handlers:       make(map[string]CommandHandler, 0),
		mut:            sync.Mutex{},
		logger:         logger,
		failedCommands: make(chan *FailedCommand),

		mutex: mutex,
	}
}

func (cb *CommandBus) RegisterCommand(command bus.Dto, handler CommandHandler) error {
	cb.mut.Lock()
	defer cb.mut.Unlock()

	commandName, err := cb.commandName(command)
	if err != nil {
		return err
	}

	if _, ok := cb.handlers[*commandName]; ok {
		return NewCommandAlreadyRegistered("command already registered", *commandName)
	}

	cb.handlers[*commandName] = handler

	return nil
}

func (cb *CommandBus) getHandler(command bus.Dto) (CommandHandler, error) {
	commandName, err := cb.commandName(command)
	if err != nil {
		return nil, err
	}
	if handler, ok := cb.handlers[*commandName]; ok {
		return handler, nil
	}

	return nil, NewCommandNotRegistered("command not registered", *commandName)
}

func (cb *CommandBus) Exec(ctx context.Context, command bus.Dto) error {
	handler, err := cb.getHandler(command)
	if err != nil {
		return err
	}

	return cb.doHandle(ctx, handler, command)
}

func (cb *CommandBus) ExecAsync(ctx context.Context, command bus.Dto) error {
	commandName, err := cb.commandName(command)
	if err != nil {
		return err
	}

	if handler, ok := cb.handlers[*commandName]; ok {
		go cb.doHandleAsync(ctx, handler, command)

		return nil
	}

	return NewCommandNotRegistered("command not registered", *commandName)
}

func (cb *CommandBus) doHandle(ctx context.Context, handler CommandHandler, command bus.Dto) error {
	if blockerCommand, ok := command.(bus.BlockOperationCommand); ok {
		operation := func() (any, error) {
			return nil, handler.Handle(ctx, blockerCommand)
		}

		_, err := cb.mutex.Mutex(ctx, blockerCommand.BlockingKey(), operation)
		if err != nil {
			cb.logger.Error(
				ctx,
				"error handling command with mutex",
				slog.String("error", err.Error()),
				slog.String("command", blockerCommand.Type()),
				slog.String("blocking_key", blockerCommand.BlockingKey()),
			)
		}

		return err
	}

	return handler.Handle(ctx, command)
}

func (cb *CommandBus) doHandleAsync(ctx context.Context, handler CommandHandler, command bus.Dto) {
	err := cb.doHandle(ctx, handler, command)

	if err != nil {
		cb.failedCommands <- &FailedCommand{
			command:        command,
			handler:        handler,
			timesProcessed: 1,
		}
		cb.logger.Error(ctx, err.Error())
	}
}

func (cb *CommandBus) commandName(cmd any) (*string, error) {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return nil, CommandNotValid{"only pointer to commands are allowed"}
	}

	name := value.String()

	return &name, nil
}

// ReprocessFailedAsyncCommands will process all failed async commands in the failedCommands channel
func (cb *CommandBus) ReprocessFailedAsyncCommands(ctx context.Context, maxTimes int) {
	for {
		select {
		case <-ctx.Done():
			close(cb.failedCommands)
			cb.logger.Warn(ctx, "exiting safely failed commands consumer...")
			return
		case failedCommand := <-cb.failedCommands:
			if failedCommand.timesProcessed >= maxTimes {
				continue
			}

			failedCommand.timesProcessed++
			if err := cb.doHandle(ctx, failedCommand.handler, failedCommand.command); err != nil {
				cb.logger.Warn(ctx, err.Error(), slog.String("previous_error", err.Error()))
			}
		}
	}
}

type CommandNotValid struct {
	message string
}

func (i CommandNotValid) Error() string {
	return i.message
}
