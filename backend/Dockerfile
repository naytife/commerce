# Use an official Go image to build the application
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY . .

# Build the Go application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

# Use scratch as base image
FROM scratch

# Copy CA certificates from the builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the binary and docs
COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/docs /app/docs

# Expose the port
EXPOSE 8000

# Run the binary
ENTRYPOINT ["/app/bin/api"]
