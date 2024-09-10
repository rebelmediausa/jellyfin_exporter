ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/jellyfin_exporter /bin/jellyfin_exporter

EXPOSE      9594
USER        nobody
ENTRYPOINT  [ "/bin/jellyfin_exporter" ]
