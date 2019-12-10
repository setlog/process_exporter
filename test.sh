#!/bin/bash
set -e
set -x
go test -v ./...
go build -o testing/integration_test github.com/setlog/process_exporter/integration_test
cd testing
./integration_test
cd ..
rm -rf testing
