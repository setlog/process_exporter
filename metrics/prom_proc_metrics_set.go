package metrics

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type ProcessMetricsSet struct {
	processMetrics map[int]*PrometheusProcessMetrics
	namespace      string
	procBinaryName string
	procArgName    string
}

func NewProcessMetricsSet(namespace, procBinaryName, procArgName string) *ProcessMetricsSet {
	return &ProcessMetricsSet{
		processMetrics: make(map[int]*PrometheusProcessMetrics),
		namespace:      namespace,
		procBinaryName: procBinaryName,
		procArgName:    procArgName,
	}
}

func (set *ProcessMetricsSet) UpdateMonitoredSet() {
	if set.processMetrics == nil {
		panic("called update on disposed ProcessMetricsSet")
	}
	processIds := findProcesses(set.procBinaryName)
	errs := AdjustMetricsMap(set.processMetrics, processIds, set.namespace)
	for _, err := range errs {
		log.Warn(fmt.Sprintf("Could not report metrics for process: %v.", err))
	}
}

func (set *ProcessMetricsSet) Dispose() {
	for _, metrics := range set.processMetrics {
		metrics.Unregister()
	}
	set.processMetrics = nil
}

func findProcesses(processName string) (pids map[int]ps.Process) {
	pids = make(map[int]ps.Process)
	procs, err := ps.Processes()
	if err != nil {
		panic(err)
	}
	for _, proc := range procs {
		name := proc.Executable()
		fmt.Println(name)
		if processName == "" || filepath.Base(name) == processName {
			pid := proc.Pid()
			pids[pid] = proc
		}
	}
	return pids
}

func AdjustMetricsMap(metricMap map[int]*PrometheusProcessMetrics, pids map[int]ps.Process, metricNamespace string) (errs []error) {
	removePids, newPids := FindPidDifferences(metricMap, pids)
	for _, pid := range newPids {
		name, err := procDescriptiveName(pid)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get descriptive process name for PID %d: %w", pid, err))
			continue
		}
		m := newProcessMetrics(pids[pid], name, metricNamespace)
		err = m.Register()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to register process metrics for PID %d: %w", pid, err))
			continue
		}
		metricMap[pid] = m
		log.Infof("Started monitoring of process \"%s\" with PID %d.", pids[pid].Executable(), pids[pid].Pid())
	}
	for _, pid := range removePids {
		m := metricMap[pid]
		m.Unregister()
		delete(metricMap, pid)
		log.Infof("Stopped monitoring of process with PID %d.", pids[pid].Pid())
	}
	return errs
}

func FindPidDifferences(pidMap map[int]*PrometheusProcessMetrics, wantedPids map[int]ps.Process) (removePids, newPids []int) {
	for pid := range wantedPids {
		if _, ok := pidMap[pid]; !ok {
			newPids = append(newPids, pid)
		}
	}
	for pid := range pidMap {
		if _, ok := wantedPids[pid]; !ok {
			removePids = append(removePids, pid)
		}
	}
	return removePids, newPids
}

func (set *ProcessMetricsSet) UpdateMetrics() {
	for pid, processMetrics := range set.processMetrics {
		updateMetrics(processMetrics, pid)
	}
}

func updateMetrics(pm *PrometheusProcessMetrics, withPid int) {
	processMetrics, err := getProcMetrics(withPid)
	if err != nil {
		log.Warn(fmt.Sprintf("Could not determine metrics for process %d: %v.", withPid, err))
		return
	}
	pm.Set(processMetrics)
}

func procDescriptiveName(pid int) (string, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}
	args, err := proc.CmdlineSlice()
	if err != nil {
		return "", err
	}
	return descriptiveNameFromArgs(args)
}

func descriptiveNameFromArgs(args []string) (string, error) {
	if len(args) <= 1 {
		return "", fmt.Errorf("too few arguments")
	}
	args = args[1:]
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.Usage = func() {}
	name := flagSet.String("name", "", "")
	err := flagSet.Parse(args)
	if err != nil {
		return "", err
	}
	if *name == "" {
		return "", fmt.Errorf("name not found or empty")
	}
	return *name, nil
}
