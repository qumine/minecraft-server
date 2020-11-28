FROM adoptopenjdk/openjdk11:ubi-minimal-jre

EXPOSE 80
EXPOSE 25565

COPY  qumine-server-java /usr/local/bin/qumine-server
RUN chmod +x /usr/local/bin/qumine-server

WORKDIR /qumine/
ENTRYPOINT [ "qumine-server", "server"]