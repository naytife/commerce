# Use an official Go image to build the application
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY . .

# Build the Go application
RUN go build -o bin/api ./cmd/api

# Use a minimal image for running the application
FROM ubuntu:jammy

# Install CA certificates
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean

# Set the working directory for the runtime
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/bin/api /app/bin/api

# Expose the port your application listens on (adjust if necessary)
EXPOSE 8000

# Run the binary
CMD ["/app/bin/api"]
