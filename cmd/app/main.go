package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/cmd/di"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/observability"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	go observability.StartMetricsServer(ctx, di.CommonServices.Logger, ":9100")

	// Create custom metrics and send data every 5 seconds
	addCounterMetricValue()
	addGaugeMetricValue()

	// Shutdown servers on SIGINT, SIGTERM or error
	select {
	case err := <-errorsChannel:
		di.ErrorShutdown(ctx, cancel, err)
	case <-ctx.Done():
		di.GracefulShutdown(ctx)
	}
}

func addCounterMetricValue() {
	opsProcessed := promauto.NewCounter(prometheus.CounterOpts{
		Name:        "cqrs_monitored_app_example_metric_total",
		Help:        "Example custom counter metric",
		ConstLabels: prometheus.Labels{"example_label": "example_label_value"},
	})
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(5 * time.Second)
		}
	}()
}

func addGaugeMetricValue() {
	opsProcessed := promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "cqrs_monitored_app_example_gauge_metric",
		Help:        "Example custom gauge metric",
		ConstLabels: prometheus.Labels{"example_label": "example_label_value"},
	})
	go func() {
		for {
			opsProcessed.Set(float64(rand.Intn(100)))
			time.Sleep(5 * time.Second)
		}
	}()
}
