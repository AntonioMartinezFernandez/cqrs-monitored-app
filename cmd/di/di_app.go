package di

import (
	"context"
	"log/slog"
)

type AppDi struct {
	CommonServices *CommonServices
	HttpServices   *HttpServices
	SystemServices *SystemModuleServices
	BookServices   *BookModuleServices
}

func InitAppDi(ctx context.Context) *AppDi {
	commonServices := InitCommonServices(ctx)
	httpServices := InitHttpServices(commonServices)
	systemServices := InitSystemModuleServices(commonServices, httpServices)
	bookServices := InitBookModuleServices(commonServices, httpServices)

	return &AppDi{
		CommonServices: commonServices,
		HttpServices:   httpServices,
		SystemServices: systemServices,
		BookServices:   bookServices,
	}
}

func (ad *AppDi) ErrorShutdown(ctx context.Context, cancel context.CancelFunc, err error) {
	defer cancel()
	if err == nil {
		return
	}

	ad.HttpServices.Router.Shutdown(ctx)

	ad.CommonServices.Logger.Error(
		ctx,
		"error on starting server",
		slog.String("service", ad.CommonServices.Config.AppServiceName),
		slog.String("version", ad.CommonServices.Config.AppVersion),
		slog.String("error", err.Error()),
	)
}

func (ad *AppDi) GracefulShutdown(ctx context.Context) {
	ad.HttpServices.Router.Shutdown(ctx)

	ad.CommonServices.Logger.Info(
		ctx,
		"server stopped",
		slog.String("service", ad.CommonServices.Config.AppServiceName),
		slog.String("version", ad.CommonServices.Config.AppVersion),
	)
}
