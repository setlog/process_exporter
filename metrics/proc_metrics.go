package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

type ProcessMetrics struct {
	cpu             float64
	ram             uint64
	swap            uint64
	diskReadBytes   uint64
	diskWriteBytes  uint64
	diskReadCount   uint64
	diskWriteCount  uint64
	networkInBytes  uint64
	networkOutBytes uint64
}

func getProcMetrics(pid int) (processMetrics *ProcessMetrics, err error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	m := &ProcessMetrics{}
	m.cpu, err = proc.CPUPercent()
	if err != nil {
		return nil, err
	}
	mem, err := proc.MemoryInfo()
	if err != nil {
		return nil, err
	}
	m.ram = mem.RSS
	m.swap = mem.Swap
	ioDisk, err := proc.IOCounters()
	if err != nil {
		return nil, err
	}
	m.diskReadBytes = ioDisk.ReadBytes
	m.diskWriteBytes = ioDisk.WriteBytes
	m.diskReadCount = ioDisk.ReadCount
	m.diskWriteCount = ioDisk.WriteCount
	ioNet, err := proc.NetIOCounters(false)
	if err != nil {
		return nil, err
	}
	if len(ioNet) > 1 {
		return nil, fmt.Errorf("got IO info for all NICs seperately when sum was requested")
	}
	if len(ioNet) == 1 {
		m.networkInBytes = ioNet[0].BytesRecv
		m.networkOutBytes = ioNet[0].BytesSent
	}
	return m, nil
}
