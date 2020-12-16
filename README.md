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

## Quick Start

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=VANILLA qumine/qumine-server-java:latest
```

## Accessing the server console

You can access the server console by executing ```console``` inside of the container. This will stream logs and allow you to send commands to the server.

# Configuration

## eula.txt

To use the server you will need to accept the eula of mojang.
```
EULA=true
```

## server.properties

Server properties can be set via the encironment variables prefixed with ```SERVER_PROPERTIES_*```.

For example to confiugre the max-players property you would set the environemt like this:
```
SERVER_PROPERTIES_MAX_PLAYERS=10
```

## whitelist.json

***whitelist.json will only be populated on the first start***

If the whitelist is turned on via ```SERVER_PROPERTIES_WHITE_LIST=true``` you can add initial players to the whitelist via the ```SERVER_WHITE_LIST``` environment variable.

For example to add the users ```Notch``` and ```Herobrine``` to the initial whitelist you would set:
```
SERVER_WHITE_LIST=Notch,Herobrine
```

You can force the whitelist.json to be overriden every time the container starts up by settings the ```SERVER_WHITE_LIST_OVERRIDE``` environment variable.
```
SERVER_WHITE_LIST_OVERRIDE=true
```

## ops.json

***ops.json will only be populated on the first start***

You can add initial players to the operators via the ```SERVER_OPS``` environment variable.

For example to add the users ```Notch``` and ```Herobrine``` to the initial operators you would set:
```
SERVER_OPS=Notch,Herobrine
```

You can force the ops.json to be overriden every time the container starts up by settings the ```SERVER_OPS_OVERRIDE``` environment variable.
```
SERVER_OPS_OVERRIDE=true
```

## Server Types

### Vanilla

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=VANILLA -e SERVER_VERSION=1.16.4 qumine/qumine-server-java:latest
```

### PaperMC

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=PAPERMC -e SERVER_VERSION=latest qumine/qumine-server-java:latest
```

### Waterfall

```
docker run -it --rm -p 8080:8080 -p 25565:25577 -e SERVER_TYPE=WATERFALL -e SERVER_VERSION=latest qumine/qumine-server-java:latest
```

### Travertine

```
docker run -it --rm -p 8080:8080 -p 25565:25577 -e SERVER_TYPE=TRAVERTINE -e SERVER_VERSION=latest qumine/qumine-server-java:latest
```

### Yatopia

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=YATOPIA -e SERVER_VERSION=latest qumine/qumine-server-java:latest
```

### Custom

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=CUSTOM -e SERVER_CUSTOM_URL=https://papermc.io/api/v1/paper/1.16.4/296/download qumine/qumine-server-java:latest
```

## Plugins

```
W.I.P
```

# Deployment

```
TODO
```

## Docker

```
TODO
```

## Kubernetes

```
TODO
```

## Helm

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