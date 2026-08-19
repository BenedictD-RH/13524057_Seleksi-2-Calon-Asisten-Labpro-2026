FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sync-worker .

FROM alpine:latest

WORKDIR /

COPY --from=builder /sync-worker /sync-worker

# Set Gin to production mode
ENV GIN_MODE=release
EXPOSE 8084

ENTRYPOINT ["/sync-worker"]