package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/event"
)

func main() {
	// Initialize Dependencies
	rootCtx, rootCancel := di.RootContext()
	defer rootCancel()
	commonServices := di.InitCommonServices(rootCtx)

	ctx, cancel := context.WithTimeout(rootCtx, 4*time.Second)
	defer cancel()

	commonServices.Logger.Info(ctx, "example started")

	// Register the commands
	commonServices.CommandBus.RegisterCommand(&BlockingCommand1{}, NewBlockingCommand1Handler(commonServices.EventBus))
	commonServices.CommandBus.RegisterCommand(&BlockingCommand2{}, NewBlockingCommand2Handler(commonServices.EventBus))

	// Create commands
	command1 := NewBlockingCommand1("command-id-1")
	command2 := NewBlockingCommand2("command-id-2")

	// Execute the command
	err := commonServices.CommandBus.Exec(ctx, command1)
	if err != nil {
		commonServices.Logger.Error(ctx, err.Error())
	}

	err = commonServices.CommandBus.Exec(ctx, command2)
	if err != nil {
		commonServices.Logger.Error(ctx, err.Error())
	}

	// Wait some time to see the output
	<-time.After(10 * time.Second)
	commonServices.Logger.Info(
		ctx,
		"example finished",
	)
}

func NewBlockingCommand1(
	id string,
) *BlockingCommand1 {
	return &BlockingCommand1{
		ID: id,
	}
}

type BlockingCommand1 struct {
	ID string
}

func (cdp *BlockingCommand1) Type() string {
	return "blocking-command-1"
}
func (cdp *BlockingCommand1) BlockingKey() string {
	return "blocking-commands"
}

func NewBlockingCommand1Handler(
	eventBus *event.EventBus,
) *BlockingCommand1Handler {
	return &BlockingCommand1Handler{
		eventBus: eventBus,
	}
}

type BlockingCommand1Handler struct {
	eventBus *event.EventBus
}

func (cb BlockingCommand1Handler) Handle(ctx context.Context, command bus.Dto) error {
	cmd, ok := command.(*BlockingCommand1)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	fmt.Println("Started handling of BlockingCommand1 with ID: ", cmd.ID)

	// Create a channel to receive the process result
	result := make(chan any)

	// Start the process execution in a goroutine
	go HeavyExecutionProcess(ctx, cmd.ID, result)

	// Wait for either the process to complete or the context to be canceled
	select {
	case <-ctx.Done():
		// Cancel the process
		close(result)
		return fmt.Errorf("process canceled: %w", ctx.Err())
	case res := <-result:
		// Process completed successfully
		fmt.Println(res.(string))
		return nil
	}
}

func NewBlockingCommand2(
	id string,
) *BlockingCommand2 {
	return &BlockingCommand2{
		ID: id,
	}
}

type BlockingCommand2 struct {
	ID string
}

func (cdp *BlockingCommand2) Type() string {
	return "blocking-command-2"
}
func (cdp *BlockingCommand2) BlockingKey() string {
	return "blocking-commands"
}

func NewBlockingCommand2Handler(
	eventBus *event.EventBus,
) *BlockingCommand2Handler {
	return &BlockingCommand2Handler{
		eventBus: eventBus,
	}
}

type BlockingCommand2Handler struct {
	eventBus *event.EventBus
}

func (cb BlockingCommand2Handler) Handle(ctx context.Context, command bus.Dto) error {
	cmd, ok := command.(*BlockingCommand2)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	fmt.Println("Started handling of BlockingCommand2 with ID: ", cmd.ID)

	// Create a channel to receive the process result
	result := make(chan any)
	defer close(result)

	// Start the process execution in a goroutine
	go HeavyExecutionProcess(ctx, cmd.ID, result)

	// Wait for either the process to complete or the context to be canceled
	select {
	case <-ctx.Done():
		// Cancel the process
		return fmt.Errorf("process canceled: %w", ctx.Err())
	case res := <-result:
		// Process completed successfully
		fmt.Println(res.(string))
		return nil
	}
}

func HeavyExecutionProcess(ctx context.Context, id string, res chan<- any) {
mainloop:
	for {
		select {
		case <-ctx.Done():
			break mainloop
		case <-time.After(3000 * time.Millisecond):
			res <- fmt.Sprintf("Finished handling of process with ID: %s", id)
			break mainloop
		}
	}
}
