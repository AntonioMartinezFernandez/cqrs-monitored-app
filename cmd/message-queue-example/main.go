package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/queue"
)

func main() {
	// Initialize Dependencies
	rootCtx, rootCancel := di.RootContext()
	defer rootCancel()
	commonServices := di.InitCommonServices(rootCtx)

	ctx, cancel := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel()

	commonServices.Logger.Info(rootCtx, "example started")

	// Subscribe handlers to the message queue
	commonServices.OldInMemoryMQ.RegisterHandler(&MyRandomMessage{}, NewMyRandomMessageHandler())
	commonServices.OldInMemoryMQ.RegisterHandler(&MyFailingMessage{}, NewMyFailingMessageHandler())
	// commonServices.InMemoryMessageQueue.DeregisterHandler(&MyRandomMessage{}, NewMyRandomMessageHandler())

	// Create messages
	message1 := NewMyRandomMessage(commonServices.UlidProvider.New().String())
	failingMessage1 := NewMyFailingMessage(commonServices.UlidProvider.New().String())
	message2 := NewMyRandomMessage(commonServices.UlidProvider.New().String())

	// Start message queue processing
	commonServices.OldInMemoryMQ.Start(ctx)

	// Publish messages
	err1 := commonServices.OldInMemoryMQ.Publish(message1)
	err2 := commonServices.OldInMemoryMQ.Publish(failingMessage1)
	err3 := commonServices.OldInMemoryMQ.Publish(message2)

	if err1 != nil || err2 != nil || err3 != nil {
		panic("error publishing messages to the queue")
	}

	// Wait a few seconds to see the output
	<-time.After(10 * time.Second)
	commonServices.Logger.Info(
		rootCtx,
		"example finished",
	)
}

// ------------------ Messages and Handlers ------------------ //
type MyRandomMessage struct {
	id          string
	messageType string
	payload     map[string]any

	maxRetries int
	maxTimeout time.Duration
}

func NewMyRandomMessage(id string) *MyRandomMessage {
	return &MyRandomMessage{
		id:          id,
		messageType: "my-random-message",
		payload:     map[string]any{"payload_key": "payload_value"},

		maxRetries: 3,
		maxTimeout: 1 * time.Second,
	}
}

func (m *MyRandomMessage) ID() string {
	return m.id
}

func (m *MyRandomMessage) MessageType() string {
	return m.messageType
}

func (m *MyRandomMessage) Payload() map[string]any {
	return m.payload
}

func (m *MyRandomMessage) MaxRetries() int {
	return m.maxRetries
}

func (m *MyRandomMessage) MaxTimeout() time.Duration {
	return m.maxTimeout
}

type MyRandomMessageHandler struct {
}

func NewMyRandomMessageHandler() *MyRandomMessageHandler {
	return &MyRandomMessageHandler{}
}

func (h *MyRandomMessageHandler) Name() string {
	return "my-random-message-handler"
}

func (h *MyRandomMessageHandler) Handle(ctx context.Context, message queue.Message) error {
	myRandomMessage, ok := message.(*MyRandomMessage)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	fmt.Println("Handling message ID", myRandomMessage.ID())
	time.Sleep(3 * time.Second)
	fmt.Println("Message ID", myRandomMessage.ID(), "handled")

	return nil
}

type MyFailingMessage struct {
	id          string
	messageType string
	payload     map[string]any

	maxRetries int
	maxTimeout time.Duration
}

func NewMyFailingMessage(id string) *MyFailingMessage {
	return &MyFailingMessage{
		id:          id,
		messageType: "my-failing-message",
		payload:     map[string]any{"payload_key": "payload_value"},

		maxRetries: 3,
		maxTimeout: 1 * time.Second,
	}
}

func (m *MyFailingMessage) ID() string {
	return m.id
}

func (m *MyFailingMessage) MessageType() string {
	return m.messageType
}

func (m *MyFailingMessage) Payload() map[string]any {
	return m.payload
}

func (m *MyFailingMessage) MaxRetries() int {
	return m.maxRetries
}

func (m *MyFailingMessage) MaxTimeout() time.Duration {
	return m.maxTimeout
}

type MyFailingMessageHandler struct {
}

func NewMyFailingMessageHandler() *MyFailingMessageHandler {
	return &MyFailingMessageHandler{}
}

func (h *MyFailingMessageHandler) Name() string {
	return "my-random-message-handler"
}

func (h *MyFailingMessageHandler) Handle(ctx context.Context, message queue.Message) error {
	myFailingMessage, ok := message.(*MyFailingMessage)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	fmt.Println("Handling message ID", myFailingMessage.ID())
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Message ID", myFailingMessage.ID(), "handled")

	return fmt.Errorf("error because this message should fail")
}
