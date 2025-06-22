# Multi-stage Dockerfile for store-deployer service
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git ca-certificates

# Copy go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -a -installsuffix cgo \
    -o store-deployer main.go

# Final runtime image - using distroless for minimal size
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary
COPY --from=builder /app/store-deployer .

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run as non-root user
USER nonroot:nonroot

# Expose port
EXPOSE 9003

# Run the service
ENTRYPOINT ["./store-deployer"]
