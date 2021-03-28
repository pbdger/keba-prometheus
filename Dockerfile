FROM alpine:3.12.4

# the nonprivileged user to start entrypoint with (will be replaced with a random userid at runtime)
ENV RUNTIMEUSER=1001
ENV TZ Europe/Berlin

ENV wallboxName localhost
ENV wallboxPort 502

EXPOSE 8080

USER root

COPY ./bin/keba-prometheus.linux /keba-prometheus
RUN chmod +x keba-prometheus

USER ${RUNTIMEUSER}

ENTRYPOINT ["/keba-prometheus"]