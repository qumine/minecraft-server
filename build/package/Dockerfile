ARG ARCH=
FROM ${ARCH}openjdk:11-jre-slim

EXPOSE 8080 25565

COPY qumine-server-java /usr/local/bin/qumine-server

RUN groupadd -rg 4242 qumine \
&& useradd -NMrg 4242 qumine \
&& mkdir /data \
&& chown -R qumine:qumine /data \
&& chmod 755 /usr/local/bin/qumine-server

USER qumine
VOLUME /data
WORKDIR /data/
ENTRYPOINT [ "qumine-server", "server"]