package system_application

import (
	"context"

	system_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/system/domain"

	pkg_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

type GetHealthcheckQueryHandler struct {
	serviceName   string
	ulidProvider  pkg_utils.UlidProvider
	healthChecker system_domain.HealthChecker
}

type GetHealthcheckQueryHandlerResponse struct {
	Id          string            `json:"id" jsonapi:"primary,healthcheck"`
	Status      map[string]string `json:"status" jsonapi:"attr,status"`
	ServiceName string            `json:"service_name" jsonapi:"attr,service_name"`
}

func NewGetHealthcheckQueryHandler(
	serviceName string,
	ulidProvider pkg_utils.UlidProvider,
	healthChecker system_domain.HealthChecker,
) *GetHealthcheckQueryHandler {
	return &GetHealthcheckQueryHandler{
		serviceName:   serviceName,
		ulidProvider:  ulidProvider,
		healthChecker: healthChecker,
	}
}

func (q GetHealthcheckQueryHandler) Handle(ctx context.Context, query pkg_bus.Dto) (any, error) {
	_, ok := query.(*GetHealthcheckQuery)
	if !ok {
		return nil, pkg_bus.NewInvalidDto("invalid query")
	}

	statuses, err := q.healthChecker.Check(ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcheckQueryHandlerResponse{
		Id:          q.ulidProvider.New().String(),
		Status:      statuses,
		ServiceName: q.serviceName,
	}, nil
}
