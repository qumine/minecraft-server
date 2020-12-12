FROM adoptopenjdk/openjdk11:alpine-jre

EXPOSE 8080 25565

COPY  qumine-server-java /usr/local/bin/qumine-server
RUN chmod +x /usr/local/bin/qumine-server

WORKDIR /data/
ENTRYPOINT [ "qumine-server", "server"]