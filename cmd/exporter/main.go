package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/setlog/process_exporter/cmd/exporter/flags"
	"github.com/setlog/process_exporter/cmd/exporter/metrics"
)

func main() {
	namespace, procBinaryName, nameFlag, port := flags.Parse(os.Args[1:])
	http.Handle("/metrics", newHttpHandler(metrics.NewProcessMetricsSet(namespace, procBinaryName, nameFlag)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
