# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN mkdir -p bin && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/pack-calculator ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/pack-calculator ./pack-calculator

# Copy templates and static files
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run the application
CMD ["./pack-calculator"]

