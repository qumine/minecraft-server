ARG ARCH=
FROM ${ARCH}openjdk:11-jre-slim

EXPOSE 8080 25565

COPY --chmod=755 qumine-server-java /usr/local/bin/qumine-server

RUN groupadd -rg 4242 qumine && useradd -NMrg 4242 qumine

USER qumine
WORKDIR /data/
ENTRYPOINT [ "qumine-server", "server"]