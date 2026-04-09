# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/app ./cmd/app


FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
# binary
COPY --from=builder /app/app /app/app

# migrations for postgres.Migrate("internal/storage/migrations", ...)
COPY --from=builder /app/internal/storage/migrations /app/internal/storage/migrations

EXPOSE 50051
CMD ["./app"]
