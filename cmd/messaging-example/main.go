package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/messaging"
)

func main() {
	// Initialize Dependencies
	rootCtx, rootCancel := di.RootContext()
	defer rootCancel()
	commonServices := di.InitCommonServices(rootCtx)

	ctx, cancel := context.WithTimeout(rootCtx, 30*time.Second)
	defer cancel()

	commonServices.Logger.Info(rootCtx, "example started")

	// Register message queue retriever
	messageRetriever := messaging.NewInMemoryMessageRetriever(commonServices.Logger, commonServices.InMemoryMessageQueue)

	// Create message consumer (with middlewares)
	middlewares := []messaging.Middleware{
		messaging.PanicRecovererMiddleware(),
	}
	consumer := messaging.NewMessageConsumer(commonServices.Logger, messageRetriever, middlewares)

	// Register message handlers
	exampleMessageHandler := NewExampleMessageHandler(commonServices.Logger)
	consumer.RegisterHandler(exampleMessageHandler)

	failingMessageHandler := NewFailingMessageHandler(commonServices.Logger)
	consumer.RegisterHandler(failingMessageHandler)

	// Create messages
	msg1, _ := NewExampleMessage(
		commonServices.UlidProvider.New().String(),
		"Michael Jordan",
		95,
		map[string]any{"foo": "bar"},
		time.Now(),
	)
	msg2, _ := NewFailingMessage(
		commonServices.UlidProvider.New().String(),
		map[string]any{"foo": "bar"},
		time.Now(),
	)
	msg3, _ := NewExampleMessage(
		commonServices.UlidProvider.New().String(),
		"Kobe Bryant",
		93,
		map[string]any{"foo": "bar"},
		time.Now(),
	)

	// Start message queue processing
	go func() {
		consumer.Run(ctx)
	}()

	// Publish messages
	err1 := commonServices.InMemoryMessageQueue.Publish(msg1)
	err2 := commonServices.InMemoryMessageQueue.Publish(msg2)
	err3 := commonServices.InMemoryMessageQueue.Publish(msg3)

	fmt.Println(commonServices.InMemoryMessageQueue.PendingMessages())

	if err1 != nil || err2 != nil || err3 != nil {
		println(err1, err2, err3)
		panic("error publishing messages to the queue")
	}

	// Wait a few seconds to see the output
	<-time.After(10 * time.Second)
	commonServices.Logger.Info(
		rootCtx,
		"example finished",
	)
}

// ------------------ Message ------------------ //
const ExampleMessageType = "company.example-entity.example-message"
const ExampleMessageSchema = "schemas/example-message.schema.json"

type ExampleMessage struct {
	*messaging.BaseMessage

	playerName  string
	playerScore int
}

func (m *ExampleMessage) Name() string {
	return m.playerName
}

func (m *ExampleMessage) Age() int {
	return m.playerScore
}

func NewExampleMessage(
	id string,
	playerName string,
	playerScore int,
	msgMetadata map[string]any,
	msgCreatedAt time.Time,
) (*ExampleMessage, error) {
	msgData := map[string]any{
		"playerName":  playerName,
		"playerScore": playerScore,
	}

	baseMessage, err := messaging.NewBaseMessage(
		id,
		ExampleMessageType,
		msgData,
		ExampleMessageSchema,
		msgMetadata,
		"example-entity",
		"example-entity-id",
		id,
		msgCreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ExampleMessage{
		BaseMessage: baseMessage,

		playerName:  playerName,
		playerScore: playerScore,
	}, nil
}

// ------------------ Message Handler ------------------ //
type ExampleMessageHandler struct {
	logger logger.Logger
}

func NewExampleMessageHandler(
	logger logger.Logger,
) *ExampleMessageHandler {
	return &ExampleMessageHandler{
		logger: logger,
	}
}

func (ch *ExampleMessageHandler) HandlerFunc() (messaging.HandlerFunc, error) {
	return func(ctx context.Context, msg messaging.Message) error {
		castedMsg, ok := msg.(*ExampleMessage)
		if !ok {
			return fmt.Errorf("unexpected message type: %T", msg)
		}

		ch.logger.Info(ctx, "processing message", slog.String("message_id", string(castedMsg.ID())))
		time.Sleep(1 * time.Second)
		ch.logger.Info(ctx, "processed message", slog.String("message_id", string(castedMsg.ID())))

		return nil
	}, nil
}

func (ch *ExampleMessageHandler) MessageType() string {
	return ExampleMessageType
}

func (ch *ExampleMessageHandler) Validator() (messaging.ValidatorFunc, error) {
	return func(ctx context.Context, msg *messaging.BaseMessage) error {
		// Validate the message payload using the JSON schema
		fmt.Println("validating with schema ", msg.DataSchema())
		return nil
	}, nil
}

func (ch *ExampleMessageHandler) MessageResolver() messaging.MessageResolver {
	return func(msg *messaging.BaseMessage) (messaging.Message, error) {
		data := msg.Data()
		playerName, ok := data["playerName"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid player name")
		}
		playerScore, ok := data["playerScore"].(int)
		if !ok {
			return nil, fmt.Errorf("invalid player score")
		}

		return NewExampleMessage(
			msg.ID(),
			playerName,
			int(playerScore),
			msg.Metadata(),
			msg.CreatedAt(),
		)
	}
}

// ------------------ Failing Message ------------------ //
const FailingMessageType = "company.example-entity.failing-message"
const FailingMessageSchema = "schemas/failing-message.schema.json"

type FailingMessage struct {
	*messaging.BaseMessage
}

func NewFailingMessage(
	id string,
	msgMetadata map[string]any,
	msgCreatedAt time.Time,
) (*FailingMessage, error) {
	msgData := map[string]any{}

	baseMessage, err := messaging.NewBaseMessage(
		id,
		FailingMessageType,
		msgData,
		FailingMessageSchema,
		msgMetadata,
		"example-entity",
		"example-entity-id",
		id,
		msgCreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &FailingMessage{
		BaseMessage: baseMessage,
	}, nil
}

// ------------------ Message Handler ------------------ //
type FailingMessageHandler struct {
	logger logger.Logger
}

func NewFailingMessageHandler(
	logger logger.Logger,
) *FailingMessageHandler {
	return &FailingMessageHandler{
		logger: logger,
	}
}

func (ch *FailingMessageHandler) HandlerFunc() (messaging.HandlerFunc, error) {
	return func(ctx context.Context, msg messaging.Message) error {
		castedMsg, ok := msg.(*FailingMessage)
		if !ok {
			return fmt.Errorf("unexpected message type: %T", msg)
		}

		ch.logger.Info(ctx, "processing message", slog.String("message_id", string(castedMsg.ID())))
		time.Sleep(100 * time.Millisecond)
		ch.logger.Info(ctx, "processed message", slog.String("message_id", string(castedMsg.ID())))

		return fmt.Errorf("failing message should fails")
	}, nil
}

func (ch *FailingMessageHandler) MessageType() string {
	return FailingMessageType
}

func (ch *FailingMessageHandler) Validator() (messaging.ValidatorFunc, error) {
	return func(ctx context.Context, msg *messaging.BaseMessage) error {
		// Validate the message payload using the JSON schema
		fmt.Println("validating with schema ", msg.DataSchema())
		return nil
	}, nil
}

func (ch *FailingMessageHandler) MessageResolver() messaging.MessageResolver {
	return func(msg *messaging.BaseMessage) (messaging.Message, error) {
		return NewFailingMessage(
			msg.ID(),
			msg.Metadata(),
			msg.CreatedAt(),
		)
	}
}
