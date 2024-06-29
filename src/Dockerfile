ARG APP_NAME=theApp

# Build stage
FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o algebror ./main.go

# Final stage
FROM docker.io/library/alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=builder /app/algebror .
EXPOSE 18080
CMD ["./algebror"]
