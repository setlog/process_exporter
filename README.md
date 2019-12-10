## process_exporter

Makes available the following metrics for Prometheus for processes matching a given name:

- CPU usage
- RAM usage
- Swap usage
- Disk IO (bytes)
- Disk IO (count)

### Building

```sh
./build.sh
```

### Usage

```sh
./process_exporter -namespace # Prometheus metric namespace
                   -binary # Name of binary to monitor
                   -nameflag # Infer value of "name"-label from value of this command line flag of monitored process
                   -port # Port to listen on for request (default: 80)
                   -interval # Interval between metrics being refreshed (default: 10 seconds)
```

### TODO

- Network IO.
