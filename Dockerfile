# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/api

# Final stage
FROM alpine:3.18

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/service .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the service
CMD ["./service"]