# Use a minimal Go image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o main .

# Use a smaller base image for production
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install Node.js, npm and CA certificates
RUN apk add --no-cache nodejs npm ca-certificates

# Copy the compiled binary from builder
COPY --from=builder /app/main .

# Create templates directory (will be mounted or provided at runtime)
RUN mkdir -p ./templates/

# Expose ports (if required, e.g., if it's an API)
EXPOSE 8001

# Run the Go binary
CMD ["./main"]