FROM docker.io/golang:1.24-alpine AS builder

COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.21.3
USER 1000
COPY --from=builder /build/aiven-metadata-prometheus-exporter /bin/aiven-metadata-prometheus-exporter
EXPOSE 2112
CMD ["/bin/aiven-metadata-prometheus-exporter"]
