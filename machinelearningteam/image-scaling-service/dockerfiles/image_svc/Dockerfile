FROM golang:1.13 AS builder

WORKDIR /app


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

ENTRYPOINT ["/app/image_svc"]
