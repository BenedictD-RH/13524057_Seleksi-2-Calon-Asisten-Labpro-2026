FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM alpine:latest

WORKDIR /

COPY --from=builder /server /server

# Set Gin to production mode
ENV GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["/server"]