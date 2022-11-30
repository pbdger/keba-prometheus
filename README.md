# Keba-Prometheus image
Provides Prometheus Metric for a Keba Wallbox as container image.

### How to start ?
Start your container binding the external port 8080.

> docker run -d --name=kebametrics -p 8080:8080 --env wallboxName=<your wallbox ip> pbdger/keba-prometheus

Try it out.

### How to get metrics ?

Open in a browser the URL with your servername and the metric port.

> http://localhost:8080/metrics

### Additional optional environment parameters
> debug: true | false

> wallboxPort: number, default is 502 

## Grafana Integration
You find a default Grafana board here: https://grafana.com/grafana/dashboards/14121

## Core binary usage
###Prerequisites
#### Mandatory
Set the mandatory environment variable (This example is Linux based);

```
export wallboxPort=<IP or servername of your wallbox, e.g. 192.168.08.15>
```

#### Optional
Set the optional environment variables (This example is Linux based);

```
export wallboxPort=<Port on which your TCP/modbus listens. Default is 502>
```

#### Call on console
```
keba-prometheus
```



### How to build your own version ?

```
GOOS=windows GOARCH=amd64 go build -o ./bin/keba-prometheus.exe keba-prometheus.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/keba-prometheus.linux keba-prometheus.go
```
