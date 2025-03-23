package messaging

import (
	"context"
)

type MessageRetriever interface {
	StartRetrieving(ctx context.Context) (<-chan *BaseMessage, error)
}
