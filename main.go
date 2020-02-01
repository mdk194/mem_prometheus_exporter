package main

import (
	"fmt"
	"net/http"
	"sync"

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

		var wg sync.WaitGroup
		var out []proc.ProcStatus
		errChan := make(chan error)
		done := make(chan interface{})

		for _, p := range procList {
			wg.Add(1)

			// Parallel read status
			go func(pid int) {
				defer wg.Done()

				ps, err := proc.NewStatus(pid, fmt.Sprintf("/proc/%d/status", pid))
				if err != nil {
					errChan <- err
				}
				out = append(out, ps)
			}(p.PID)
		}

		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case err := <-errChan:
			return nil, fmt.Errorf("Failed to read status %v", err)
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
