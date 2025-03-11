package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

	// Create custom metric and send data every 5 seconds
	customCounterMetric, _ := di.CommonServices.Observability.Meter.Int64Counter(
		"cqrs-monitored-app.example-metric", // name of the metric
	)
	go counterAddMetric(ctx, di.CommonServices.Logger, customCounterMetric)

	// Shutdown servers on SIGINT, SIGTERM or error
	select {
	case err := <-errorsChannel:
		di.ErrorShutdown(ctx, cancel, err)
	case <-ctx.Done():
		di.GracefulShutdown(ctx)
	}
}

func counterAddMetric(ctx context.Context, l logger.Logger, customCounterMetric metric.Int64Counter) {
	for range time.Tick(5 * time.Second) {
		l.Info(
			ctx,
			"sending custom counter metric",
			slog.String(
				"metric_name",
				"cqrs_monitored_app_example_metric_total",
			),
		)

		// Record values to the metric instruments and add labels
		customCounterMetric.Add(ctx, 1, metric.WithAttributes(
			attribute.String("custom_string_attribute", "custom_string_attribute_value"),
			attribute.Int("custom_int_attribute", 1),
		))
	}
}
