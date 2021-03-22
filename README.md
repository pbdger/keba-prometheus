# keba-prometheus
Provides data as Prometheus Metric for a Kuba Wallbox.

Usage
-----
Before usage:
Set the mandatory environment variable (This example is Linux based);

```go
export wallboxPort=<IP or servername of your wallbox, e.g. 192.168.08.15>
```
Set the optional environment variables (This example is Linux based);

```go
export wallboxPort=<Port on which your TCP/modbus listens. Default is 502>
```

Basic usage:
```go
keba-prometheus
```

How to get metrics ?:

Open in a browser the URL with your servername and the metric port.

```go
http://localhost:8080/metrics
```

How to build your own version ?

```go
GOOS=windows GOARCH=amd64 go build -o ./bin/keba-prometheus.exe keba-prometheus.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/keba-prometheus.linux keba-prometheus.go
```
