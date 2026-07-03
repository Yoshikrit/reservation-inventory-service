FROM golang:1.26.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api  ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/grpc ./cmd/grpc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/consumer ./cmd/consumer

FROM gcr.io/distroless/static-debian13:nonroot

USER 65532:65532

WORKDIR /app

COPY --from=builder /app/bin/api  .
COPY --from=builder /app/bin/grpc .
COPY --from=builder /app/bin/consumer .
