# Backend Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o /goapp main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Установка git в конечном образе
RUN apk update && apk add --no-cache git

COPY --from=builder /goapp .

# Command to run the executable
CMD ["./goapp"]
