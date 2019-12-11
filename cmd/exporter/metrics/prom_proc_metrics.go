package metrics

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-ps"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type PrometheusProcessMetrics struct {
	registeredCollectors      []prometheus.Collector
	previousMetrics           *ProcessMetrics
	cpuGauge                  prometheus.Gauge
	ramGauge                  prometheus.Gauge
	swapGauge                 prometheus.Gauge
	storageReadBytesGauge     prometheus.Gauge
	storageWriteBytesGauge    prometheus.Gauge
	storagediskReadCountGauge prometheus.Gauge
	storageWriteCountGauge    prometheus.Gauge
}

func newPrometheusProcessMetrics(proc ps.Process, descriptiveName, metricNamespace string) (processMetrics *PrometheusProcessMetrics) {
	processMetrics = &PrometheusProcessMetrics{}
	binaryName := filepath.Base(proc.Executable())
	pid := fmt.Sprintf("%d", proc.Pid())
	processMetrics.makeGauges(metricNamespace, pid, binaryName, descriptiveName)
	processMetrics.makeStorageGauges(metricNamespace, pid, binaryName, descriptiveName)
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

func (pm *PrometheusProcessMetrics) makeStorageGauges(metricNamespace, pid, binaryName, descriptiveName string) {
	pm.storageReadBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "storage_read_bytes",
		Help:        "Total read from storage (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.storageWriteBytesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "storage_write_bytes",
		Help:        "Total written to storage (bytes)",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.storagediskReadCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "storage_reads",
		Help:        "Total reads from storage",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
	pm.storageWriteCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   metricNamespace,
		Name:        "storage_writes",
		Help:        "Total writes to storage",
		ConstLabels: prometheus.Labels{"pid": pid, "bin": binaryName, "name": descriptiveName},
	})
}

func (pm *PrometheusProcessMetrics) Register() error {
	pm.Unregister()
	var err error
	for _, collector := range []prometheus.Collector{
		pm.cpuGauge,
		pm.ramGauge,
		pm.swapGauge,
		pm.storageReadBytesGauge,
		pm.storageWriteBytesGauge,
		pm.storagediskReadCountGauge,
		pm.storageWriteCountGauge,
	} {
		err = prometheus.Register(collector)
		if err != nil {
			break
		}
		pm.registeredCollectors = append(pm.registeredCollectors, collector)
	}
	if err != nil {
		pm.Unregister()
	}
	return err
}

func (pm *PrometheusProcessMetrics) Unregister() {
	for _, collector := range pm.registeredCollectors {
		prometheus.Unregister(collector)
	}
	pm.registeredCollectors = nil
}

func (pm *PrometheusProcessMetrics) Set(processMetrics *ProcessMetrics) {
	if pm.previousMetrics != nil {
		deltaTime := processMetrics.cpuSampleTime.Sub(pm.previousMetrics.cpuSampleTime).Seconds()
		if deltaTime > 0 {
			pm.cpuGauge.Set(((processMetrics.cpuDuration - pm.previousMetrics.cpuDuration) * 100) / deltaTime)
		} else {
			log.Warn("deltaTime <= 0")
		}
	}
	pm.ramGauge.Set(float64(processMetrics.ram))
	pm.swapGauge.Set(float64(processMetrics.swap))
	pm.storageReadBytesGauge.Set(float64(processMetrics.storageReadBytes))
	pm.storageWriteBytesGauge.Set(float64(processMetrics.storageWriteBytes))
	pm.storagediskReadCountGauge.Set(float64(processMetrics.storageReadCount))
	pm.storageWriteCountGauge.Set(float64(processMetrics.storageWriteCount))
	pm.previousMetrics = processMetrics
}
