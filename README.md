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
- [ ] Basic API health endpoints.
- [ ] Basic download of server plugins.
- [ ] Basic updating of server plugins.
- [ ] GRPC API for controlling the server remotely(start, stop, ).
- [ ] GRPC API for log streaming.
- [ ] GRPC API for console streaming(in, out).
- [ ] Cryo, stop server if no client is connected.

# Usage

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