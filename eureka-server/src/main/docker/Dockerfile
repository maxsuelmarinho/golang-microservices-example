FROM openjdk:8-alpine

ENV LANG en_US.utf8
ENV TZ "America/Sao_Paulo"

VOLUME /tmp
VOLUME /app/logs
VOLUME /app/config

RUN apk update && \
    apk add --no-cache bash gettext curl

ADD eureka-server-*.jar /app/app.jar
RUN sh -c 'touch /app/app.jar'
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

HEALTHCHECK --interval=10s --timeout=10s CMD [ "curl", "--silent", "--fail", "localhost:8761/actuator/health" ] || exit 1

EXPOSE 8761
ENTRYPOINT ["/docker-entrypoint.sh"]