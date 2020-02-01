package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	http.Handle("/metrics", promhttp.Handler())
	logger.Info("Starting exporter at :8080/metrics")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Failed to start", zap.Error(err))
	}
}
