FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app-b-backend .

FROM alpine:latest

WORKDIR /

COPY --from=builder /app-b-backend /app-b-backend


# Set Gin to production mode
ENV GIN_MODE=release
EXPOSE 8791

ENTRYPOINT ["/app-b-backend"]