# Sentry Gateway

Sentry Gateway is an HTTP reverse proxy that can monitor and rate limit the requests to your self-hosted Sentry.
If you have experience downtimes because of burst events from your product to your self-hosted Sentry, Sentry Gateway can be useful.
It can be built into a single binary file or a docker image.

## Deploy

### 1- Binary file

#### Dependencies

- Go1.14
- Git

#### Steps

- Open a terminal and Clone the project by the command `git clone git@github.com:cafebazaar/sentry-gateway.git` 
- `cd sentry-gateway`
- Build the source code: `go build ./cmd -o sentrygatewayd`
- Copy `config.yaml` file and customize it
- Export `SENTRY_GATEWAY_CONFIG_FILE_PATH` variable to customized config file. For example: `export SENTRY_GATEWAY_CONFIG_FILE_PATH=/home/user/my-sentry-gateway-config.yaml`
- Run project `./sentrygatewayd`

### 2- Dockerfile

#### Steps

- You can build docker image `docker build . -t sentry-gateway`. Or pull docker image from docker hub.
`docker pull bardia13/sentry-gateway:latest`
- Write your own config file.
- Run docker image: `docker run -e SENTRY_GATEWAY_CONFIG_FILE_PATH=/config.yaml -v config.yaml:/config.yaml bardia13/sentry-gateway`

## Config file

Default config file is in root of the project 

Variable | Description | Preferred value
---------|-------------|----------------
proxy.targetAddress | Address of your private Sentry | 127.0.0.1:9000
proxy.timeout | Timeout is the maximum amount of time a dial will wait for a connect to complete. If Deadline is also set, it may fail earlier. The default is no timeout. When using TCP and dialing a host name with multiple IP addresses, the timeout may be divided between them. With or without a timeout, the operating system may impose its own earlier timeout. For instance, TCP timeouts are often around 3 minutes. | 25s
proxy.keepAlive | KeepAlive specifies the interval between keep-alive probes for an active network connection. If zero, keep-alive probes are sent with a default value (currently 15 seconds), if supported by the protocol and operating system. Network protocols or operating systems that do not support keep-alives ignore this field. If negative, keep-alive probes are disabled. | 20s
proxy.maxIdleConns |  MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit. | 20
proxy.idleConnTimeout | IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself. Zero means no limit. | 5s
proxy.expectContinueTimeout | ExpectContinueTimeout, if non-zero, specifies the amount of time to wait for a server's first response headers after fully  writing the request headers if the request has an "Expect: 100-continue" header. Zero means no timeout and causes the body to be sent immediately, without waiting for the server to approve. This time does not include the time to send the request header. | 2s
throttle.RPS | Request Per Second limit | 20
throttle.burst | Maximum number of request in one second | 30
metrics.enabled: | Prometheus metrics is enabled | true
metrics.listenAddress | The interface that prometheus is listening to | 0.0.0.0:9090
logLevel | Log level (panic, fatal, error, warn, info, debug, trace) | info


## Contribution and Future works

All contributions and PRs are welcomed.

Some of the features that will be helpful for others:

- TLS support
- Load balancer
