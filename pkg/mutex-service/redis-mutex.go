package distributed_sync

import (
	"context"
	"errors"
	"log/slog"
	"time"

	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisMutexService struct {
	sync   *redsync.Redsync
	logger pkg_logger.Logger
}

func NewRedisMutexService(redisClient *redis.Client, logger pkg_logger.Logger) *RedisMutexService {
	pool := goredis.NewPool(redisClient)

	return &RedisMutexService{sync: redsync.New(pool), logger: logger}
}

func (rm *RedisMutexService) Mutex(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	mutex := rm.sync.NewMutex(
		mutexName+":"+key,
		redsync.WithExpiry(30*time.Second),
		redsync.WithRetryDelay(25*time.Millisecond),
		redsync.WithTimeoutFactor(0.05),
	)

	if _, err := utils.RetryFunc(func() (interface{}, error) {
		return nil, rm.adquireLock(ctx, mutex)
	}, 4); err != nil {
		rm.logger.Error(ctx, "error locking mutex sync", slog.String("err", err.Error()), slog.String("mutex_key", mutex.Name()))
		return nil, NewErrorLockMutexKey(key, err)
	}

	result, err := fn()
	if _, err := utils.RetryFunc(func() (interface{}, error) {
		return nil, rm.releaseLock(ctx, mutex)
	}, 4); err != nil {
		rm.logger.Error(ctx, "error unlocking mutex sync", slog.String("error", err.Error()), slog.String("mutex_key", mutex.Name()))
		return nil, NewErrorReleaseLockMutexKey(key, err)
	}

	return result, err
}

func (rm *RedisMutexService) releaseLock(ctx context.Context, mutex *redsync.Mutex) error {
	if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
		rm.logger.Warn(ctx, "error unlocking mutex sync - retrying", slog.String("error", err.Error()), slog.String("mutex_key", mutex.Name()))
		if err != nil {
			return err
		}

		return errors.New("redis mutex invalid status when unlocking")
	}

	return nil
}

func (rm *RedisMutexService) adquireLock(ctx context.Context, mutex *redsync.Mutex) error {

	if err := mutex.LockContext(ctx); err != nil {
		rm.logger.Warn(ctx, "error locking mutex sync - retrying", slog.String("error", err.Error()), slog.String("mutex_key", mutex.Name()))
		return err
	}

	return nil
}
