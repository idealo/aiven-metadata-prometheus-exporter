FROM golang:1.18-alpine AS builder

ENV AIVEN_API_TOKEN;

COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.16.0
COPY --from=builder /build/aiven-metadata-prometheus-exporter /bin/aiven-metadata-prometheus-exporter
EXPOSE 2112
CMD ["/bin/aiven-metadata-prometheus-exporter"]
