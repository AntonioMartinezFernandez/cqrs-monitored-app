package distributed_sync

import (
	"context"
	"sync"

	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
)

type InmemoryMutexService struct {
	logger pkg_logger.Logger
	mut    sync.Mutex
	locks  map[string]*sync.Mutex
}

func NewInmemoryMutexService(
	logger pkg_logger.Logger,
) *InmemoryMutexService {
	return &InmemoryMutexService{
		logger: logger,
		mut:    sync.Mutex{},
		locks:  make(map[string]*sync.Mutex),
	}
}

func (im *InmemoryMutexService) Mutex(ctx context.Context, key string, fn func() (any, error)) (any, error) {
	im.mut.Lock()
	keyMut, ok := im.locks[key]
	if !ok {
		keyMut = &sync.Mutex{}
		im.locks[key] = keyMut
	}
	im.mut.Unlock()

	// Try to acquire the lock with a timeout
	done := make(chan struct{})
	var result any
	var err error

	keyMut.Lock()
	defer keyMut.Unlock()
	go func() {
		result, err = fn()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return result, err
	}
}
