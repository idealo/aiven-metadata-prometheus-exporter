FROM docker.io/golang:1.23-alpine AS builder

COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.20.3
USER 1000
COPY --from=builder /build/aiven-metadata-prometheus-exporter /bin/aiven-metadata-prometheus-exporter
EXPOSE 2112
CMD ["/bin/aiven-metadata-prometheus-exporter"]
