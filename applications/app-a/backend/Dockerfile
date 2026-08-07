FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app-a-backend .

FROM alpine:latest

WORKDIR /

COPY --from=builder /app-a-backend /app-a-backend


# Set Gin to production mode
ENV GIN_MODE=release
EXPOSE 8691

ENTRYPOINT ["/app-a-backend"]