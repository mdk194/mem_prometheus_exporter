package main

import (
	"fmt"

	"github.com/mdk194/mem_prometheus_exporter/proc"
)

func main() {
	procs, err := proc.AllProcs()
	if err != nil {
		fmt.Println("Error listing all current procs", err)
	}
	for _, p := range procs {
		fmt.Printf("%d ", p.PID)
	}
	fmt.Println()
}
