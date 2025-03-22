package observability

import (
	"context"
	"net/http"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer starts a HTTP server to serve Prometheus metrics.
// Example usage: StartMetricsServer(ctx, logger, ":9100")
func StartMetricsServer(ctx context.Context, l logger.Logger, listenAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		l.Error(ctx, "error starting metrics server", logger.ErrValue("error", err))
	}
}
