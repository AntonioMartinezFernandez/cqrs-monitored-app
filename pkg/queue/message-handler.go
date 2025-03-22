package queue

import "context"

type MessageHandler interface {
	Name() string
	Handle(ctx context.Context, message Message) error
}
