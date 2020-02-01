package main

import (
	"strconv"

	"github.com/mdk194/mem_prometheus_exporter/proc"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &collector{}

type collector struct {
	ProcessMemory *prometheus.Desc
	stats         func() ([]proc.ProcStatus, error)
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
