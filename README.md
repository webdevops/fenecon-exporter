# Fenecon Prometheus exporter

[![license](https://img.shields.io/github/license/webdevops/fenecon-exporter.svg)](https://github.com/webdevops/fenecon-exporter/blob/master/LICENSE)
[![DockerHub](https://img.shields.io/badge/DockerHub-webdevops%2Ffenecon--exporter-blue)](https://hub.docker.com/r/webdevops/fenecon-exporter/)
[![Quay.io](https://img.shields.io/badge/Quay.io-webdevops%2Ffenecon--exporter-blue)](https://quay.io/repository/webdevops/fenecon-exporter)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/fenecon-exporter)](https://artifacthub.io/packages/search?repo=fenecon-exporter)

Prometheus exporter for Fenecon systems

## Usage

```
Usage:
  fenecon-exporter [OPTIONS]

Application Options:
      --log.debug                    debug mode [$LOG_DEBUG]
      --log.devel                    development mode [$LOG_DEVEL]
      --log.json                     Switch log output to json format [$LOG_JSON]
      --fenecon.request.timeout=     Request timeout (default: 10s) [$FENECON_REQUEST_TIMEOUT]
      --fenecon.request.parallel=    Number of parallel requests (default: 1) [$FENECON_REQUEST_PARALLEL]
      --fenecon.request.retries=     Request retries (default: 1) [$FENECON_REQUEST_RETRIES]
      --fenecon.request.waittime=    Request retries (default: 2s) [$FENECON_REQUEST_WAITTIME]
      --fenecon.request.maxwaittime= Request retries (default: 5s) [$FENECON_REQUEST_MAXWAITTIME]
      --fenecon.auth.username=       Username for fenecon login [$FENECON_AUTH_USERNAME]
      --fenecon.auth.password=       Password for fenecon login (default: user) [$FENECON_AUTH_PASSWORD]
      --server.bind=                 Server address (default: :8080) [$SERVER_BIND]
      --server.timeout.read=         Server read timeout (default: 5s) [$SERVER_TIMEOUT_READ]
      --server.timeout.write=        Server write timeout (default: 60s) [$SERVER_TIMEOUT_WRITE]

Help Options:
  -h, --help                         Show this help message
```

## HTTP Endpoints

| Endpoint       | Description                         |
|----------------|-------------------------------------|
| `/metrics`     | Default prometheus golang metrics   |
| `/probe`       | Probe metrics from Fenecon system   |

### /probe/metrics parameters

request metrics from Fenecon system

| GET parameter | Default | Required | Type                    | Description                                |
|---------------|---------|----------|-------------------------|--------------------------------------------|
| `target`      |         | **yes**  | string                  | Url to Fenecon system, eg `http://fenecon` |
| `ess`         | `0`     | no       | integer (positive only) | Number of batteries (`ess`) modules        |
| `charger`     | `0`     | no       | integer (positive only) | Number of chargers (`charger`) modules     |
| `meter`       | `0`     | no       | integer (positive only) | Number of meters (`meter`) modules         |
