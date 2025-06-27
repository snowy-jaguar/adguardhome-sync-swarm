FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app

RUN apk update && apk add upx ca-certificates tzdata

ARG VERSION=main
ARG BUILD="N/A"

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux

COPY . /go/src/app/

RUN go build -a -installsuffix cgo -ldflags="-w -s -X github.com/bakito/adguardhome-sync/version.Version=${VERSION} -X github.com/bakito/adguardhome-sync/version.Build=${BUILD}" -o adguardhome-sync .

RUN go version && upx -q adguardhome-sync

# Final image (using alpine to allow entrypoint scripts)
FROM alpine:3.18

LABEL maintainer="snowy-jaguar <github@snowy-jaguar.com>"


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
