QuMine - Server - Java
---
![GitHub Release](https://img.shields.io/github/v/release/qumine/minecraft-server)
![GitHub Workflow](https://img.shields.io/github/workflow/status/qumine/minecraft-server/release)
[![GoDoc](https://godoc.org/github.com/qumine/minecraft-server?status.svg)](https://godoc.org/github.com/qumine/minecraft-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/qumine/minecraft-server)](https://goreportcard.com/report/github.com/qumine/minecraft-server)

Docker Image for running minecraft servers.

# Status

- [X] Basic download of server JAR.
- [X] Basic updating of server JAR.
- [X] Basic wrapping of JVM process.
- [X] Basic API health endpoints.
- [X] Basic download of server plugins.
- [ ] Basic updating of server plugins.
- [X] GRPC API for controlling the server remotely(start, stop, ).
- [X] GRPC API for log streaming.
- [X] GRPC API for console streaming(in, out).
- [ ] Cryo, stop server if no client is connected.

# Usage

## Quick Start

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=VANILLA qumine/minecraft-server:latest
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

## AIKAR flags

You can use AIKAR's optimized flags for better server performance with certain use cases.

```
USE_AIKAR_FLAGS=true
```

## Additional files

To automatically download and extract additional files you can use the ADDITIONAL_FILES environment variable.

***Currently only HTTP and HTTP/S sources are supported***
***Currently only ZIP sources are supported***

```
ADDITIONAL_FILES="https://example.com/file.zip,https://example.com/file2.zip"
```

## Plugins

To automatically download plugins and keep them updated you can use the SERVER_PLUGINS environment variable.

***Currently only HTTP and HTTP/S plugins sources are supported***
```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=PAPERMC -e SERVER_VERSION=latest -e SERVER_PLUGINS="https://github.com/BlueMap-Minecraft/BlueMap/releases/download/v1.3.0-snap/BlueMap-1.3.0-snap-spigot.jar,https://ci.opencollab.dev/job/GeyserMC/job/Floodgate/job/master/lastSuccessfulBuild/artifact/bukkit/target/floodgate-bukkit.jar" qumine/minecraft-server:latest
```

## Server Types

### Custom

In custom mode by default the provided url will be downloaded and later executed with the java -jar command.
```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=CUSTOM -e SERVER_CUSTOM_URL=https://papermc.io/api/v1/paper/1.16.4/296/download qumine/minecraft-server:latest
```

If you need to customize the startup command of your server you can use the SERVER_CUSTOM_COMMAND and SERVER_CUSTOM_ARGS environment variables.
```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=CUSTOM -e SERVER_CUSTOM_URL=https://papermc.io/api/v1/paper/1.16.4/296/download -e SERVER_CUSTOM_COMMAND=java -e SERVER_CUSTOM_ARGS="-XX:+UseG1GC,-jar,download,nogui" qumine/minecraft-server:latest
```

### PaperMC

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=PAPERMC -e SERVER_VERSION=latest qumine/minecraft-server:latest
```

### ServerStarter

This will configure the server using [ServerStarter](https://github.com/Yoosk/ServerStarter).
```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=PAPERMC -e SERVER_VERSION=latest qumine/minecraft-server:latest
```

**NOTE:** The amount of memory must fit the minRam and maxRam options of the server-setup-config.yaml.

### Travertine

```
docker run -it --rm -p 8080:8080 -p 25565:25577 -e SERVER_TYPE=TRAVERTINE -e SERVER_VERSION=latest qumine/minecraft-server:latest
```

### Vanilla

```
docker run -it --rm -p 8080:8080 -p 25565:25565 -e EULA=true -e SERVER_TYPE=VANILLA -e SERVER_VERSION=1.16.4 qumine/minecraft-server:latest
```

### Waterfall

```
docker run -it --rm -p 8080:8080 -p 25565:25577 -e SERVER_TYPE=WATERFALL -e SERVER_VERSION=latest qumine/minecraft-server:latest
```

# Deployment

## Operator

**W.I.P**: An operator based approach for managing Minecraft Servers inside of kubernetes will follow in the future.

## Docker

```
docker run -it --rm -p 8080:8080 -p 25565:25565 \
  -e EULA=true \
  -e SERVER_TYPE=VANILLA \
  -e SERVER_VERSION=1.16.4 \
  -e SERVER_WHITE_LIST=User1,User2 \
  -e SERVER_OPS=User1,User2 \
  -e SERVER_PROPERTIES_MOTD="Example Minecraft Server" \
  qumine/minecraft-server:latest
```

## Kubernetes

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-minecraft-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: minecraft-server
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: example-minecraft-server
    spec:
      containers:
      - env:
        - name: EULA
          value: "true"
        - name: SERVER_TYPE
          value: VANILLA
        - name: SERVER_VERSION
          value: 1.16.4
        - name: SERVER_WHITE_LIST
          value: User1,User2
        - name: SERVER_OPS
          value: User1,User2
        - name: SERVER_PROPERTIES_MOTD
          value: Example Minecraft Server
        image: docker.io/qumine/minecraft-server:v0.1.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 5
          httpGet:
            path: /health/live
            port: 8080
            scheme: HTTP
          periodSeconds: 5
          successThreshold: 1
          timeoutSeconds: 1
        name: minecraft-server
        ports:
        - containerPort: 25565
          name: minecraft
          protocol: TCP
        - containerPort: 8080
          name: metrics
          protocol: TCP
        readinessProbe:
          failureThreshold: 5
          httpGet:
            path: /health/ready
            port: 8080
            scheme: HTTP
          periodSeconds: 5
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            cpu: "2"
            memory: 4000Mi
          requests:
            cpu: "2"
            memory: 4000Mi
        startupProbe:
          failureThreshold: 24
          httpGet:
            path: /health/ready
            port: 8080
            scheme: HTTP
          periodSeconds: 5
          successThreshold: 1
          timeoutSeconds: 1
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