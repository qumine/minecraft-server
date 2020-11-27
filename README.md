QuMine - Server - Java
---
![GitHub Release](https://img.shields.io/github/v/release/qumine/qumine-server-java)
![GitHub Workflow](https://img.shields.io/github/workflow/status/qumine/qumine-server-java/release)
[![GoDoc](https://godoc.org/github.com/qumine/qumine-server-java?status.svg)](https://godoc.org/github.com/qumine/qumine-server-java)
[![Go Report Card](https://goreportcard.com/badge/github.com/qumine/qumine-server-java)](https://goreportcard.com/report/github.com/qumine/qumine-server-java)

Docker Image for running minecraft servers.

# Status

- [X] Basic download of server JAR.
- [X] Basic updating of server JAR.
- [X] Basic wrapping of JVM process.
- [X] Basic API health endpoints.
- [ ] Basic download of server plugins.
- [ ] Basic updating of server plugins.
- [X] GRPC API for controlling the server remotely(start, stop, ).
- [X] GRPC API for log streaming.
- [X] GRPC API for console streaming(in, out).
- [ ] Cryo, stop server if no client is connected.

# Usage

## Configuration

Configuration is done via environment variables.

## Accessing the server console

You can access the server console by executing ```/qumine-server console``` inside of the container. This will stream logs and allow you to send commands to the server.

## Kubernetes

*HELM Charts can be found here: [qumine/charts](https://github.com/qumine/charts)*

# Development

## Perfrom a Snapshot release locally

```
docker run -it --rm \
  -v ${PWD}:/build -w /build \
  -v /var/run/docker.sock:/var/run/docker.sock \
  goreleaser/goreleaser \
  release --snapshot --rm-dist
```