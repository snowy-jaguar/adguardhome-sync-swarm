# Dockerfile for adguardhome-sync
# This Dockerfile builds a Go application and packages it into a minimal Alpine Linux image.
# The application is a sync tool for AdGuard Home, which allows syncing DNS records and other
# configurations with an AdGuard Home instance.

##### Building the application #####
FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app

RUN apk update && apk add upx ca-certificates tzdata

ARG VERSION=main
ARG BUILD="N/A"

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux

COPY . /go/src/app/

ARG VERSION=main
ARG BUILD="N/A"

RUN go build -a -installsuffix cgo -ldflags="-w -s \
  -X github.com/snowy-jaguar/adguardhome-sync-swarm/version.Version=${VERSION} \
  -X github.com/snowy-jaguar/adguardhome-sync-swarm/version.Build=${BUILD}" \
  -o adguardhome-sync .


RUN go version && upx -q adguardhome-sync

##### Creating the final image #####
FROM alpine:3.18

LABEL maintainer="snowy-jaguar <contact@snowyjaguar.xyz>" \
      org.opencontainers.image.authors="snowy-jaguar"


RUN apk add --no-cache ca-certificates curl tzdata

WORKDIR /opt/go

COPY --from=builder /go/src/app/adguardhome-sync /opt/go/adguardhome-sync
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY healthcheck.sh /usr/local/bin/healthcheck.sh

RUN chmod +x /usr/local/bin/entrypoint.sh /usr/local/bin/healthcheck.sh

EXPOSE 8080

USER 1001

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run", "--config", "/config/adguardhome-sync.yaml"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 CMD /usr/local/bin/healthcheck.sh
