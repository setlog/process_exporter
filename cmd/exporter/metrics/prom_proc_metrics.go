package metrics

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type PrometheusProcessMetrics struct {
	previousMetrics      *ProcessMetrics
	cpuGauge             prometheus.Gauge
	ramGauge             prometheus.Gauge
	swapGauge            prometheus.Gauge
	diskReadBytesGauge   prometheus.Gauge
	diskWriteBytesGauge  prometheus.Gauge
	diskReadCountGauge   prometheus.Gauge
	diskWriteCountGauge  prometheus.Gauge
	networkInBytesGauge  prometheus.Gauge
	networkOutBytesGauge prometheus.Gauge
}

func newPrometheusProcessMetrics(proc ps.Process, descriptiveName, metricNamespace string) (processMetrics *PrometheusProcessMetrics) {
	processMetrics = &PrometheusProcessMetrics{}
	binaryName := filepath.Base(proc.Executable())
	pid := fmt.Sprintf("%d", proc.Pid())
	processMetrics.makeGauges(metricNamespace, pid, binaryName, descriptiveName)
	processMetrics.makeDiskGauges(metricNamespace, pid, binaryName, descriptiveName)
	processMetrics.makeNetworkGauges(metricNamespace, pid, binaryName, descriptiveName)
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

func (pm *PrometheusProcessMetrics) makeDiskGauges(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.diskReadBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "disk_read_bytes",
		Help:        "Total read from disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskWriteBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "disk_write_bytes",
		Help:        "Total written to disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskReadCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "disk_reads",
		Help:        "Total reads from disk",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.diskWriteCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "disk_writes",
		Help:        "Total writes to disk",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
}

func (pm *PrometheusProcessMetrics) makeNetworkGauges(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.networkInBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "net_read_bytes",
		Help:        "Total read from disk (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.networkOutBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
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
		pm.diskReadBytesGauge,
		pm.diskWriteBytesGauge,
		pm.diskReadCountGauge,
		pm.diskWriteCountGauge,
		pm.networkInBytesGauge,
		pm.networkOutBytesGauge,
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
	if pm.previousMetrics != nil {
		deltaTime := processMetrics.cpuSampleTime.Sub(pm.previousMetrics.cpuSampleTime).Seconds()
		if deltaTime > 0 {
			pm.cpuGauge.Set((processMetrics.cpuDuration - pm.previousMetrics.cpuDuration) / deltaTime)
		} else {
			log.Warn("deltaTime <= 0")
		}
	}
	pm.ramGauge.Set(float64(processMetrics.ram))
	pm.swapGauge.Set(float64(processMetrics.swap))
	pm.diskReadBytesGauge.Set(float64(processMetrics.diskReadBytes))
	pm.diskWriteBytesGauge.Set(float64(processMetrics.diskWriteBytes))
	pm.diskReadCountGauge.Set(float64(processMetrics.diskReadCount))
	pm.diskWriteCountGauge.Set(float64(processMetrics.diskWriteCount))
	pm.networkInBytesGauge.Set(float64(processMetrics.networkInBytes))
	pm.networkOutBytesGauge.Set(float64(processMetrics.networkOutBytes))
	pm.previousMetrics = processMetrics
}
