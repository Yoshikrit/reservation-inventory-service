FROM golang:1.26.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api  ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/grpc ./cmd/grpc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/consumer ./cmd/consumer

FROM alpine:3.23.4

RUN apk add --no-cache tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/api  .
COPY --from=builder /app/bin/grpc .
COPY --from=builder /app/bin/consumer .

USER appuser
