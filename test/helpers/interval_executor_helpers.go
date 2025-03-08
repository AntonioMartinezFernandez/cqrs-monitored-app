package helpers

import (
	"context"
	"sync"
	"time"

	pkg_logger "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

type IntervalExecutorTestSuite struct {
	ctx       context.Context
	logger    pkg_logger.Logger
	executors []pkg_utils.ExecutorFunc
	wg        *sync.WaitGroup
}

func NewIntervalExecutorTestSuite(logger pkg_logger.Logger, executors []pkg_utils.ExecutorFunc) *IntervalExecutorTestSuite {
	return &IntervalExecutorTestSuite{
		ctx:       context.Background(),
		logger:    logger,
		executors: executors,
		wg:        &sync.WaitGroup{},
	}
}

func (suite *IntervalExecutorTestSuite) Start() {
	suite.wg.Add(len(suite.executors))

	for _, executor := range suite.executors {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
		pkg_utils.IntervalExecutor(ctx, executor, suite.logger, time.NewTicker(60*time.Second), suite.wg)
	}

	suite.wg.Wait()
}
