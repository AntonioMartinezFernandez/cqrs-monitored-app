package messaging

import (
	"context"
	"errors"
	"log/slog"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
)

// InMemoryQueueConfig is the specific configuration for the message retriever
type InMemoryQueueConfig struct {
	MaxMessages int32
}

func DefaultInMemoryQueueConfig() InMemoryQueueConfig {
	return InMemoryQueueConfig{
		MaxMessages: 1,
	}
}

// Create the specific message retriever to be injected in the consumer
type InMemoryMessageRetriever struct {
	logger       logger.Logger
	messageQueue *InMemoryQueue
}

func NewInMemoryMessageRetriever(
	logger logger.Logger,
	messageQueue *InMemoryQueue,
) *InMemoryMessageRetriever {
	return &InMemoryMessageRetriever{
		logger:       logger,
		messageQueue: messageQueue,
	}
}

func (mr *InMemoryMessageRetriever) StartRetrieving(ctx context.Context) (<-chan *BaseMessage, error) {
	messages := make(chan *BaseMessage)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			default:
				err := mr.pullMessages(ctx, messages)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						mr.logger.Warn(ctx, "error pulling messages from message queue", slog.String("error", err.Error()))
					}
					continue
				}
			}
		}
	}()

	return messages, nil
}

func (mr *InMemoryMessageRetriever) pullMessages(
	ctx context.Context,
	output chan<- *BaseMessage,
) error {
	msgs, _ := mr.messageQueue.RetrieveMessages(1)

	for _, baseMsg := range msgs {
		err := mr.deleteMessage(baseMsg.ExternalID())
		if err != nil {
			mr.logger.Warn(
				ctx,
				"error deleting message from in memory message queue",
				slog.String("error", err.Error()),
				slog.String("message_id", baseMsg.ID()),
				slog.String("baseMsg_external_id", baseMsg.ExternalID()),
			)
		}

		output <- &baseMsg
	}

	return nil
}

func (mr *InMemoryMessageRetriever) deleteMessage(messageIdentifier string) error {
	return mr.messageQueue.DeleteMessage(messageIdentifier)
}
