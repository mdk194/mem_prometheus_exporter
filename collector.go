package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/mdk194/mem_prometheus_exporter/proc"
)

var _ prometheus.Collector = &collector{}

type collector struct {
	ProcessMemory *prometheus.Desc
	stats         func() ([]proc.ProcStatus, error)
}

func init() {
	c := newCollector(stats)
	prometheus.MustRegister(c)
}

func stats() ([]proc.ProcStatus, error) {
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

func newCollector(stats func() ([]proc.ProcStatus, error)) prometheus.Collector {
	return &collector{
		ProcessMemory: prometheus.NewDesc(
			"process_memory_rss_bytes",
			"Size of memory resient set size of process read from /proc/[pid]/status",
			[]string{"pid", "name"},
			nil,
		),
		stats: stats,
	}
}

// Describe implements prometheus.Collector
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.ProcessMemory,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements prometheus.Collector
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.ProcessMemory, err)
	}

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.ProcessMemory,
			prometheus.GaugeValue,
			float64(s.VmRSS),
			strconv.Itoa(s.PID), s.Name,
		)
	}
}
