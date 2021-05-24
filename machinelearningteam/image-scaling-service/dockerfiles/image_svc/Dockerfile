FROM golang:1.13 AS builder

WORKDIR /app


RUN wget -O grpc-health-probe -q https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64 \
    && chmod +x grpc-health-probe
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build the application
RUN go build cmd/image_svc/image_svc.go

FROM debian:stretch-slim AS runner

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=builder /app/image_svc /app/image_svc
COPY --from=builder /app/grpc-health-probe /usr/local/bin/grpc-health-probe

ENTRYPOINT ["/app/image_svc"]
