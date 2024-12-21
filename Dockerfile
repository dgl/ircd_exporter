ARG GO_VERSION=1.23.4
FROM golang:${GO_VERSION}-alpine AS builder

SHELL ["/bin/sh", "-euxo", "pipefail", "-c"]

ENV GOPATH=/go
ENV CGO_ENABLED=0

COPY . /go/src/ircd_exporter

WORKDIR /go/src/ircd_exporter

#RUN ls -lah ; \
#    rm go.mod go.sum ; \
#    go mod init github.com/dgl/ircd_exporter ; \
#    go mod tidy

RUN go mod download ; \
    go build -o ${GOPATH}/bin/ircd_exporter ./cmd/ircd_exporter ; \
    ${GOPATH}/bin/ircd_exporter --help


FROM scratch

ARG BUILD_VERSION
ARG BUILD_DATE
ARG BUILD_COMMIT_SHA

LABEL org.opencontainers.image.title="ircd_exporter" \
      org.opencontainers.image.version="${BUILD_VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${BUILD_COMMIT_SHA}" \
      org.opencontainers.image.description="Prometheus exporter for IRC server state" \
      org.opencontainers.image.documentation="https://github.com/dgl/ircd_exporter" \
      org.opencontainers.image.base.name="scratch" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/dgl/ircd_exporter"

COPY --from=builder --chown=65534:65534 /go/bin/ircd_exporter /usr/local/bin/ircd_exporter

# user: nobody
USER 65534

EXPOSE 9678/tcp

ENTRYPOINT ["ircd_exporter"]
#CMD ["--help"]