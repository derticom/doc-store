FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o doc-store ./cmd/doc-store

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/doc-store .
COPY config/config.yml ./config.yml
COPY --from=builder /app/migrations /app/migrations

ENTRYPOINT ["./doc-store", "-config=config.yml"]
