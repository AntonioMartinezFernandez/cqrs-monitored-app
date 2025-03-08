package di

import (
	system_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/system/application"
	system_infra "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/system/infra"
	system_http "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/system/infra/http"
)

type SystemModuleServices struct {
	HealthcheckQueryHandler *system_application.GetHealthcheckQueryHandler
}

func InitSystemModuleServices(commonServices *CommonServices, httpServices *HttpServices) *SystemModuleServices {
	healthchecker := system_infra.NewSimpleHealthChecker(commonServices.Observability.Tracer)
	healthcheckQueryHandler := system_application.NewGetHealthcheckQueryHandler(
		commonServices.Config.AppServiceName,
		commonServices.UlidProvider,
		healthchecker,
	)

	systemModuleServices := &SystemModuleServices{
		HealthcheckQueryHandler: healthcheckQueryHandler,
	}

	registerSystemQueryHandlers(commonServices, systemModuleServices)
	registerSystemRoutes(commonServices, httpServices)

	return systemModuleServices
}

func registerSystemQueryHandlers(commonServices *CommonServices, systemModuleServices *SystemModuleServices) {
	registerQueryOrPanic(
		commonServices.QueryBus,
		&system_application.GetHealthcheckQuery{},
		systemModuleServices.HealthcheckQueryHandler,
	)
}

func registerSystemRoutes(commonServices *CommonServices, httpServices *HttpServices) {
	httpServices.Router.Get(
		"/system/healthcheck",
		system_http.NewGetHealthcheckController(
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)
}
