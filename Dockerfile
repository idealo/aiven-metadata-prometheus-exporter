FROM golang:1.18-alpine AS builder

COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.16.0
COPY --from=builder /build/aiven-prometheus-exporter /bin/aiven-prometheus-exporter
EXPOSE 2112
CMD ["/bin/aiven-prometheus-exporter"]
