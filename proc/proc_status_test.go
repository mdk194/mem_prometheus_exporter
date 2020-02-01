package proc

import (
	"io/ioutil"
	"os"
	"testing"
)

var mockStatus = `Name:	systemd
Umask:	0000
State:	S (sleeping)
Tgid:	1
Ngid:	0
Pid:	1
PPid:	0
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	128
Groups:	 
NStgid:	1
NSpid:	1
NSpgid:	1
NSsid:	1
VmPeak:	  238376 kB
VmSize:	  172840 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	   10612 kB
VmRSS:	   10612 kB
RssAnon:	    2556 kB
RssFile:	    8056 kB
RssShmem:	       0 kB
VmData:	   25464 kB
VmStk:	     132 kB
VmExe:	     908 kB
VmLib:	    7996 kB
VmPTE:	      96 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
CoreDumping:	0
Threads:	1
SigQ:	0/62811
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	7be3c0fe28014a03
SigIgn:	0000000000001000
SigCgt:	00000001800004ec
CapInh:	0000000000000000
CapPrm:	0000003fffffffff
CapEff:	0000003fffffffff
CapBnd:	0000003fffffffff
CapAmb:	0000000000000000
NoNewPrivs:	0
Seccomp:	0
Speculation_Store_Bypass:	thread vulnerable
Cpus_allowed:	ff
Cpus_allowed_list:	0-7
Mems_allowed:	00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	120965
nonvoluntary_ctxt_switches:	6455`

func TestProcStatus(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "mockStatus")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(mockStatus)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	s, err := NewStatus(1, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want int
		have int
	}{
		{name: "Pid", want: 1, have: s.PID},
		{name: "Tgid", want: 1, have: s.TGID},
		{name: "VmPeak", want: 238376 * 1024, have: int(s.VmPeak)},
		{name: "VmSize", want: 172840 * 1024, have: int(s.VmSize)},
		{name: "VmLck", want: 0 * 1024, have: int(s.VmLck)},
		{name: "VmPin", want: 0 * 1024, have: int(s.VmPin)},
		{name: "VmHWM", want: 10612 * 1024, have: int(s.VmHWM)},
		{name: "VmRSS", want: 10612 * 1024, have: int(s.VmRSS)},
		{name: "RssAnon", want: 2556 * 1024, have: int(s.RssAnon)},
		{name: "RssFile", want: 8056 * 1024, have: int(s.RssFile)},
		{name: "RssShmem", want: 0 * 1024, have: int(s.RssShmem)},
		{name: "VmData", want: 25464 * 1024, have: int(s.VmData)},
		{name: "VmStk", want: 132 * 1024, have: int(s.VmStk)},
		{name: "VmExe", want: 908 * 1024, have: int(s.VmExe)},
		{name: "VmLib", want: 7996 * 1024, have: int(s.VmLib)},
		{name: "VmPTE", want: 96 * 1024, have: int(s.VmPTE)},
		{name: "VmPMD", want: 0 * 1024, have: int(s.VmPMD)},
		{name: "VmSwap", want: 0 * 1024, have: int(s.VmSwap)},
		{name: "HugetlbPages", want: 0 * 1024, have: int(s.HugetlbPages)},
		{name: "VoluntaryCtxtSwitches", want: 120965, have: int(s.VoluntaryCtxtSwitches)},
		{name: "NonVoluntaryCtxtSwitches", want: 6455, have: int(s.NonVoluntaryCtxtSwitches)},
	} {
		if test.want != test.have {
			t.Errorf("want %s %d, have %d", test.name, test.want, test.have)
		}
	}
}
