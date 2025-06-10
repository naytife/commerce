# Use an official Go image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY . .

# Generate Swagger documentation
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest
RUN swag init -g main.go -o docs -d cmd/api,internal/api/handlers,internal/api/models,internal/db --v3.1

# Build the Go application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

# Use scratch - the smallest possible base image (empty filesystem)
FROM scratch

# Copy CA certificates from builder stage for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set working directory
WORKDIR /app

# Copy the binary and docs
COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/docs /app/docs

# Expose the port
EXPOSE 8000

# Run the binary (note: no shell available in scratch, so use exec form)
ENTRYPOINT ["/app/bin/api"]
