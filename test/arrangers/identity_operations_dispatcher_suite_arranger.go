package arrangers

import (
	"context"
	"sync"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
)

type IdentityOperationsDispatcherSuiteArranger struct {
	wg *sync.WaitGroup
}

func NewIdentityOperationsDispatcherSuiteArranger(
	common *di.CommonServices,
) *IdentityOperationsDispatcherSuiteArranger {
	return &IdentityOperationsDispatcherSuiteArranger{
		wg: &sync.WaitGroup{},
	}
}

func (sa *IdentityOperationsDispatcherSuiteArranger) Arrange(ctx context.Context) {
	sa.wg.Add(0)

	sa.wg.Wait()
}
