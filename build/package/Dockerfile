ARG ARCH=
FROM ${ARCH}openjdk:11-jre-slim

VOLUME /data
EXPOSE 8080 25565

COPY entrypoint.sh /usr/local/bin/entrypoint
COPY qumine-server-java /usr/local/bin/qumine-server

RUN groupadd -g 4242 -r qumine \
&& useradd -g 4242 -M -N -r -u 4242 qumine \
&& mkdir /data \
&& chown -R qumine:qumine /data \
&& chmod -R 777 /data \
&& chmod -R 755 /usr/local/bin/entrypoint /usr/local/bin/qumine-server

USER qumine
WORKDIR /data/
ENTRYPOINT [ "entrypoint"]