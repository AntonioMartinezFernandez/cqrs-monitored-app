package distributed_sync

import (
	"context"
)

const mutexName = "mutex-service-mutex"

type MutexService interface {
	Mutex(ctx context.Context, key string, fn func() (any, error)) (any, error)
}
