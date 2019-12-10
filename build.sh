#!/bin/bash
set -e
set -x
go build -o process_exporter github.com/setlog/process_exporter/cmd/exporter
