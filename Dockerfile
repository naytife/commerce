FROM golang:1.22 AS builder

# Set the working directory to the backend
WORKDIR /app/backend

# Copy the Go files
COPY ./backend/go.* ./
COPY ./backend/*.go ./

# Download dependencies and build the binary
RUN go mod tidy
RUN go build -o bin/api

# Start a new stage for the production environment
FROM ubuntu:jammy
WORKDIR /app
COPY --from=builder /app/backend/bin/api /app/bin/api

# Run the built binary
CMD ["/app/bin/api"]
