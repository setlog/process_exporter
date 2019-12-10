package metrics_test

import (
	"reflect"
	"testing"

	"github.com/mitchellh/go-ps"
	"github.com/setlog/process_exporter/cmd/exporter/metrics"
)

type proc struct {
	pid int
}

func (p *proc) Pid() int {
	return p.pid
}

func (p *proc) PPid() int {
	return 42
}

func (p *proc) Executable() string {
	return "/home/u/bin/answertoeverything"
}

func TestFindPidDifferences(t *testing.T) {
	metricMap := make(map[int]*metrics.PrometheusProcessMetrics)
	metricMap[5] = nil
	metricMap[7] = nil
	metricMap[11] = nil
	metricMap[13] = nil
	processMap := make(map[int]ps.Process)
	processMap[5] = &proc{5}
	processMap[11] = &proc{11}
	processMap[17] = &proc{17}
	processMap[19] = &proc{19}
	removePids, newPids := metrics.FindPidDifferences(metricMap, processMap)
	if len(removePids) != 2 {
		t.Fatalf("len(removePids) was %d. Expected 2.", len(removePids))
	}
	if len(newPids) != 2 {
		t.Fatalf("len(newPids) was %d. Expected 2.", len(newPids))
	}
	if !reflect.DeepEqual(removePids, []int{7, 13}) && !reflect.DeepEqual(removePids, []int{13, 7}) {
		t.Fatalf("removePids was %v. Expected {7, 13}.", removePids)
	}
	if !reflect.DeepEqual(newPids, []int{17, 19}) && !reflect.DeepEqual(newPids, []int{19, 17}) {
		t.Fatalf("newPids was %v. Expected {17, 19}.", newPids)
	}
}
