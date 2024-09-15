# Use an official Go image to build the application
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app/backend

# Copy the Go module files to download dependencies first
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY ./backend ./

# Build the Go application
RUN go build -o bin/api ./cmd/api

# Use a minimal image for running the application
FROM ubuntu:jammy

# Set the working directory for the runtime
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/backend/bin/api /app/bin/api

# Expose the port your application listens on (adjust if necessary)
EXPOSE 8080

# Run the binary
CMD ["/app/bin/api"]
