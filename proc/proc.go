package proc

import (
	"fmt"
	"os"
	"strconv"
)

type Proc struct {
	PID int
}

type Procs []Proc

// AllProcs returns a list of all currently available processes.
func AllProcs() (Procs, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return Procs{}, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return Procs{}, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	p := Procs{}
	for _, n := range names {
		pid, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			continue
		}
		p = append(p, Proc{PID: int(pid)})
	}

	return p, nil
}
