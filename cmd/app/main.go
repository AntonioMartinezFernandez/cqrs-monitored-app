package main

import (
	"fmt"
	"log/slog"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
)

func main() {
	// Initialize Dependencies
	ctx, cancel := di.RootContext()
	defer cancel()
	errorsChannel := make(chan error)
	di := di.InitAppDi(ctx)

	di.CommonServices.Logger.Info(
		ctx,
		"starting HTTP server...",
		slog.String("service", di.CommonServices.Config.AppServiceName),
		slog.String("version", di.CommonServices.Config.AppVersion),
	)

	// Start Http Server
	go func() {
		errorsChannel <- di.HttpServices.Router.ListenAndServe(
			fmt.Sprintf("%s:%s", di.CommonServices.Config.HttpHost, di.CommonServices.Config.HttpPort),
		)
	}()

	// Shutdown servers on SIGINT, SIGTERM or error
	select {
	case err := <-errorsChannel:
		di.ErrorShutdown(ctx, cancel, err)
	case <-ctx.Done():
		di.GracefulShutdown(ctx)
	}
}
