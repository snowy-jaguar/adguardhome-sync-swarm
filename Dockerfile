# Copyright 2025 snowy-jaguar
# Contact: @snowyjaguar (Discord)
# Contact: contact@snowyjaguar.xyz (Email)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

RUN go build -a -installsuffix cgo -ldflags="-w -s \
  -X github.com/snowy-jaguar/adguardhomesync-swarm/version.Version=${VERSION} \
  -X github.com/snowy-jaguar/adguardhomesync-swarm/version.Build=${BUILD}" \
  -o adguardhome-sync .


RUN go version && upx -q adguardhome-sync

##### Creating the final image #####
FROM alpine:3.18

ENV PATH="/opt/go:$PATH"

LABEL maintainer="snowy-jaguar <contact@snowyjaguar.xyz>" \
      org.opencontainers.image.authors="snowy-jaguar"


RUN apk add --no-cache ca-certificates curl tzdata

WORKDIR /opt/go

COPY --from=builder /go/src/app/adguardhome-sync /opt/go/adguardhome-sync
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY healthcheck.sh /usr/local/bin/healthcheck.sh

RUN chmod +x /usr/local/bin/entrypoint.sh
# RUN chmod +x /usr/local/bin/healthcheck.sh

EXPOSE 8080

USER 1001

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run", "--config", "/config/adguardhome-sync.yaml"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 CMD /usr/local/bin/healthcheck.sh
