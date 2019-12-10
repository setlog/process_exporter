package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/setlog/process_exporter/cmd/exporter/flags"
	"github.com/setlog/process_exporter/cmd/exporter/metrics"
	log "github.com/sirupsen/logrus"
)

func main() {
	namespace, procBinaryName, nameFlag, port := flags.Parse(os.Args[1:])
	metricsSet := metrics.NewPrometheusProcessMetricsSet(namespace, procBinaryName, nameFlag)
	mu := &sync.Mutex{}
	updateMetricsSet(metricsSet, mu)
	// ctx, cancelFunc := context.WithCancel(context.Background())
	go keepMetricsUpToDate(metricsSet, mu)
	http.Handle("/metrics", newHttpMetricsRequestHandler(metricsSet, mu))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func keepMetricsUpToDate(set *metrics.PrometheusProcessMetricsSet, mu *sync.Mutex) {
	defer func() {
		if r := recover(); r != nil {
			defer func() { go func() { time.Sleep(time.Second); log.Exit(1) }() }()
			log.Panicf("Panic in keepMetricsUpToDate(): %v", r)
		}
	}()
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			{
				updateMetricsSet(set, mu)
			}
		}
	}
}

func updateMetricsSet(set *metrics.PrometheusProcessMetricsSet, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	set.UpdateMonitoredSet()
	set.UpdateMetrics()
}
