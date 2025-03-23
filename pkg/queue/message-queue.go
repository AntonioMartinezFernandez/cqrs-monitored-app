package queue

import "context"

type MessageQueue interface {
	Publish(message Message) error
	RegisterHandler(message Message, handler MessageHandler) error
	DeregisterHandler(message Message, handler MessageHandler) error
	Start(ctx context.Context) error
	EnqueuedMessages() int
}
