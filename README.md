## process_exporter

Makes available the following metrics for Prometheus for processes matching a given name:

- CPU usage
- RAM usage
- Swap usage
- Disk IO (bytes)
- Disk IO (count)

### Requirements

- Go 1.13

### Building

```sh
./build.sh
```

### Usage

```sh
./process_exporter -namespace # Prometheus metric namespace (default: "my")
                   -binary # Name of binary to monitor
                   -nameflag # Infer value of "name"-label from value of this command line flag of monitored process (default: "name")
                   -port # Port to listen on for request (default: 80)
                   -interval # Interval between metrics being refreshed, in seconds (default: 10)
```

### Metrics

Example output from `curl localhost/metrics` for `process_exporter -binary top -nameflag d` for `top -d 1`. (Exploiting top's `-d` parameter as a unique identifier)

```sh
# HELP my_cpu Process CPU usage (%)
# TYPE my_cpu gauge
my_cpu{bin="top",name="1",pid="26436"} 3.0630784711756576

# HELP my_ram Process RAM usage (bytes)
# TYPE my_ram gauge
my_ram{bin="top",name="1",pid="26436"} 3.80928e+06

# HELP my_storage_read_bytes Total read from storage (bytes)
# TYPE my_storage_read_bytes gauge
my_storage_read_bytes{bin="top",name="1",pid="26436"} 0

# HELP my_storage_reads Total reads from storage
# TYPE my_storage_reads gauge
my_storage_reads{bin="top",name="1",pid="26436"} 27973

# HELP my_storage_write_bytes Total written to storage (bytes)
# TYPE my_storage_write_bytes gauge
my_storage_write_bytes{bin="top",name="1",pid="26436"} 0

# HELP my_storage_writes Total writes to storage
# TYPE my_storage_writes gauge
my_storage_writes{bin="top",name="1",pid="26436"} 91

# HELP my_swap Process swap usage (bytes)
# TYPE my_swap gauge
my_swap{bin="top",name="1",pid="26436"} 0
```

### TODO

- Network IO.
