package metrics

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusProcessMetrics struct {
	lastMetricsSnapshot    *ProcessMetrics
	cpuGauge               prometheus.Gauge
	ramGauge               prometheus.Gauge
	swapGauge              prometheus.Gauge
	diskReadBytesCounter   prometheus.Counter
	diskWriteBytesCounter  prometheus.Counter
	diskReadCountCounter   prometheus.Counter
	diskWriteCountCounter  prometheus.Counter
	networkInBytesCounter  prometheus.Counter
	networkOutBytesCounter prometheus.Counter
}

func newProcessMetrics(proc ps.Process, descriptiveName, metricNamespace string) (processMetrics *PrometheusProcessMetrics) {
	processMetrics = &PrometheusProcessMetrics{}
	processMetrics.lastMetricsSnapshot = &ProcessMetrics{}
	binaryName := filepath.Base(proc.Executable())
	pid := fmt.Sprintf("%d", proc.Pid())
	processMetrics.makeGauges(metricNamespace, pid, binaryName, descriptiveName)
	processMetrics.makeDiskCounters(metricNamespace, pid, binaryName, descriptiveName)
	processMetrics.makeNetworkCounters(metricNamespace, pid, binaryName, descriptiveName)
	return processMetrics
}

func (pm *PrometheusProcessMetrics) makeGauges(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.cpuGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "cpu",
		Help:        "Process CPU usage (%)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.ramGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "ram",
		Help:        "Process RAM usage (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.swapGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "swap",
		Help:        "Process swap usage (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
}

func (pm *PrometheusProcessMetrics) makeDiskCounters(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.diskReadBytesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "disk_read_bytes",
		Help:        "Total read from disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskWriteBytesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "disk_write_bytes",
		Help:        "Total written to disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskReadCountCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "disk_reads",
		Help:        "Total reads from disk",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskWriteCountCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "disk_writes",
		Help:        "Total writes to disk",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
}

func (pm *PrometheusProcessMetrics) makeNetworkCounters(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.networkInBytesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "net_read_bytes",
		Help:        "Total read from disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.networkOutBytesCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   metricNamespace,
		Name:        "net_write_bytes",
		Help:        "Total written to disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
}

func (pm *PrometheusProcessMetrics) Register() error {
	registeredCollectors := make([]prometheus.Collector, 0)
	var err error
	for _, collector := range []prometheus.Collector{
		pm.cpuGauge,
		pm.ramGauge,
		pm.swapGauge,
		pm.diskReadBytesCounter,
		pm.diskWriteBytesCounter,
		pm.diskReadCountCounter,
		pm.diskWriteCountCounter,
		pm.networkInBytesCounter,
		pm.networkOutBytesCounter,
	} {
		err = prometheus.Register(collector)
		if err != nil {
			break
		}
		registeredCollectors = append(registeredCollectors, collector)
	}
	if err != nil {
		for _, collector := range registeredCollectors {
			prometheus.Unregister(collector)
		}
	}
	return err
}

func (pm *PrometheusProcessMetrics) Unregister() {
	prometheus.Unregister(pm.cpuGauge)
}

func (pm *PrometheusProcessMetrics) Update() {
	prometheus.Unregister(pm.cpuGauge)
}

func (pm *PrometheusProcessMetrics) Set(processMetrics *ProcessMetrics) {
	pm.cpuGauge.Set(processMetrics.cpu)
	pm.ramGauge.Set(float64(processMetrics.ram))
	pm.swapGauge.Set(float64(processMetrics.swap))
	pm.diskReadBytesCounter.Add(float64(processMetrics.diskReadBytes - pm.lastMetricsSnapshot.diskReadBytes))
	pm.diskWriteBytesCounter.Add(float64(processMetrics.diskWriteBytes - pm.lastMetricsSnapshot.diskWriteBytes))
	pm.diskReadCountCounter.Add(float64(processMetrics.diskReadCount - pm.lastMetricsSnapshot.diskReadCount))
	pm.diskWriteCountCounter.Add(float64(processMetrics.diskWriteCount - pm.lastMetricsSnapshot.diskWriteCount))
	pm.networkInBytesCounter.Add(float64(processMetrics.networkInBytes - pm.lastMetricsSnapshot.networkInBytes))
	pm.networkOutBytesCounter.Add(float64(processMetrics.networkOutBytes - pm.lastMetricsSnapshot.networkOutBytes))
	pm.lastMetricsSnapshot = processMetrics
}
