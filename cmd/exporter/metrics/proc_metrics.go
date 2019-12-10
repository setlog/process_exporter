package metrics

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/process"
)

type ProcessMetrics struct {
	cpuDuration    float64
	cpuSampleTime  time.Time
	ram            uint64
	swap           uint64
	diskReadBytes  uint64
	diskWriteBytes uint64
	diskReadCount  uint64
	diskWriteCount uint64
}

func getProcMetrics(pid int) (processMetrics *ProcessMetrics, err error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, fmt.Errorf("could not open PID %d: %w", pid, err)
	}
	m := &ProcessMetrics{}
	timeStat, err := proc.Times()
	if err != nil {
		return nil, fmt.Errorf("could not read CPU times of PID %d: %w", pid, err)
	}
	m.cpuDuration, m.cpuSampleTime = timeStat.Total(), time.Now()
	mem, err := proc.MemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("could not read memory info of PID %d: %w", pid, err)
	}
	m.ram = mem.RSS
	m.swap = mem.Swap
	ioDisk, err := proc.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("could not read disk IO info of PID %d: %w", pid, err)
	}
	m.diskReadBytes = ioDisk.ReadBytes
	m.diskWriteBytes = ioDisk.WriteBytes
	m.diskReadCount = ioDisk.ReadCount
	m.diskWriteCount = ioDisk.WriteCount
	return m, nil
}
