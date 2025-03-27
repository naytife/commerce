# Use a minimal Go image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Copy go modules and download dependencies
COPY cloud-build/go.mod cloud-build/go.sum ./
RUN go mod download

# Copy templates directory
COPY templates/ /app/templates/

# Copy the source code from cloud-build directory
COPY cloud-build/ .

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

# Copy templates folder if needed
COPY --from=builder /app/templates/ ./templates/

# Expose ports (if required, e.g., if it's an API)
EXPOSE 8001

# Run the Go binary
CMD ["./main"]