FROM adoptopenjdk/openjdk8:aarch64-ubi-minimal-jre8u272-b10

EXPOSE 80
EXPOSE 25565

COPY  qumine-server-java qumine-server
USER nobody

ENTRYPOINT [ "./qumine-server" ]