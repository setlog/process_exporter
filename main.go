package main

import (
	"log"
	"net/http"
	"os"

	"github.com/setlog/libvirt_exporter/flags"
	"github.com/setlog/libvirt_exporter/metrics"
)

func main() {
	namespace, procBinaryName, argName := flags.Parse(os.Args[1:])
	http.Handle("/metrics", newHttpHandler(metrics.NewProcessMetricsSet(namespace, procBinaryName, argName)))
	log.Fatal(http.ListenAndServe(":http", nil))
}
