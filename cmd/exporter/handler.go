package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/setlog/process_exporter/cmd/exporter/metrics"
)

type httpHandler struct {
	metricsHandler http.Handler
	metricsSet     *metrics.ProcessMetricsSet
}

func newHttpHandler(set *metrics.ProcessMetricsSet) *httpHandler {
	return &httpHandler{
		metricsHandler: promhttp.Handler(),
		metricsSet:     set,
	}
}

func (h *httpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.metricsSet.UpdateMonitoredSet()
	h.metricsSet.UpdateMetrics()
	h.metricsHandler.ServeHTTP(response, request)
}
