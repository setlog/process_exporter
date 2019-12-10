package main

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/setlog/process_exporter/cmd/exporter/metrics"
)

type httpMetricsRequestHandler struct {
	metricsHandler http.Handler
	metricsSet     *metrics.PrometheusProcessMetricsSet
	metricsMutex   *sync.Mutex
}

func newHttpMetricsRequestHandler(set *metrics.PrometheusProcessMetricsSet, metricsMutex *sync.Mutex) *httpMetricsRequestHandler {
	return &httpMetricsRequestHandler{
		metricsHandler: promhttp.Handler(),
		metricsSet:     set,
		metricsMutex:   metricsMutex,
	}
}

func (h *httpMetricsRequestHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.metricsMutex.Lock()
	defer h.metricsMutex.Unlock()
	h.metricsHandler.ServeHTTP(response, request)
}
