package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/mdk194/mem_prometheus_exporter/proc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	stats := func() ([]proc.ProcStatus, error) {
		procList, err := proc.AllProcs()
		if err != nil {
			return nil, fmt.Errorf("Failed to list processes %v", err)
		}

		var out []proc.ProcStatus
		for _, p := range procList {
			ps, err := proc.NewStatus(p.PID, fmt.Sprintf("/proc/%d/status", p.PID))
			if err != nil {
				return nil, fmt.Errorf("Failed to read /proc/%d/status %v", p.PID, err)
			}
			out = append(out, ps)
		}
		return out, nil
	}

	c := newCollector(stats)
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	logger.Info("Starting exporter at :8080/metrics")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Failed to start", zap.Error(err))
	}
}
