package di

import (
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/configs"

	pkg_http_server "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/http-server"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
	pkg_observability "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/observability"
)

type HttpServices struct {
	Router                    *pkg_http_server.Router
	JsonApiResponseMiddleware *pkg_json_api.JsonApiResponseMiddleware
}

func InitHttpServices(commonServices *CommonServices) *HttpServices {
	return &HttpServices{
		Router:                    newRouter(commonServices.Config, commonServices),
		JsonApiResponseMiddleware: pkg_json_api.NewJsonApiResponseMiddleware(commonServices.Logger),
	}
}

func newRouter(config configs.Config, commonServices *CommonServices) *pkg_http_server.Router {
	return pkg_http_server.DefaultRouter(
		config.HttpWriteTimeout,
		config.HttpReadTimeout,
		pkg_http_server.NewPanicRecoverMiddleware(commonServices.Logger).Middleware,
		pkg_observability.NewOtelInstrumentationMiddleware(commonServices.Config.AppServiceName).Middleware,
	)
}
