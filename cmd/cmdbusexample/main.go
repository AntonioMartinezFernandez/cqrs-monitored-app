package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

func main() {
	// Initialize Dependencies
	ctx, cancel := di.RootContext()
	di := di.InitAppDi(ctx)
	defer func() {
		cancel()
	}()

	di.CommonServices.Logger.Info(
		ctx,
		"starting app...",
		slog.String("service", di.CommonServices.Config.AppServiceName),
		slog.String("version", di.CommonServices.Config.AppVersion),
	)

	di.CommonServices.CommandBus.RegisterCommand(&ExampleCommand{}, NewExampleCommandHandler())

	go di.CommonServices.CommandBus.Dispatch(ctx, &ExampleCommand{Name: "example_command"})
	go di.CommonServices.CommandBus.Dispatch(ctx, &ExampleCommand{Name: "example_command"})
	// go di.CommonServices.CommandBus.Dispatch(ctx, &ExampleCommand{Name: "example_command_2"})
	// go di.CommonServices.CommandBus.Dispatch(ctx, &ExampleCommand{Name: "example_command_2"})
	// go di.CommonServices.CommandBus.Dispatch(ctx, &ExampleCommand{Name: "example_command_3"})

	<-time.After(10 * time.Second)
}

const exampleCmdName = "example_command"

type ExampleCommand struct {
	Name string
}

func (cdp *ExampleCommand) Type() string {
	return exampleCmdName
}

func (cdp *ExampleCommand) BlockingKey() string {
	return cdp.Name
}

type ExampleCommandHandler struct{}

func NewExampleCommandHandler() ExampleCommandHandler {
	return ExampleCommandHandler{}
}

func (fd ExampleCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	dpCommand, ok := command.(*ExampleCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	fmt.Printf("executing command '%s'...\n", dpCommand.Name)
	time.Sleep(3 * time.Second)
	fmt.Printf("execution of command '%s' finished\n", dpCommand.Name)

	return nil
}
