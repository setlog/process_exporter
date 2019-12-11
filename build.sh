#!/bin/bash
set -e
set -x
export GO111MODULE=on
go build -o process_exporter github.com/setlog/process_exporter/cmd/exporter
