package main

import (
	"fmt"

	"github.com/mdk194/mem_prometheus_exporter/proc"
)

func main() {
	ps, err := proc.NewStatus(1, "/proc/1/status")
	if err != nil {
		fmt.Println("Error reading proc status", err)
	}
	fmt.Println("peak RSS, RSS", ps.VmHWM, ps.VmRSS)
}
