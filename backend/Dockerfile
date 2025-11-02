# Use an official Go image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY . .

# Install Atlas for database migrations
RUN curl -sSf https://atlasgo.sh | sh -s -- --yes --version v0.35.0

# Verify Atlas installation location
RUN which atlas && atlas version

# Generate Swagger documentation
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest
RUN swag init -g main.go -o docs -d cmd/api,internal/api/handlers,internal/api/models,internal/db --v3.1

# Build the Go application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

# Use Alpine Linux for shell support and minimal footprint
FROM alpine:3.19

# Install CA certificates and create app user
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary, docs, atlas binary, and migration files
COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/docs /app/docs
COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas
COPY --from=builder /app/atlas.hcl /app/atlas.hcl
COPY --from=builder /app/internal/db/migrations /app/migrations

# Change ownership to app user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the port
EXPOSE 8000

# Create entrypoint script that runs migrations then starts the app
COPY --from=builder --chown=appuser:appgroup /app/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Use the entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]
