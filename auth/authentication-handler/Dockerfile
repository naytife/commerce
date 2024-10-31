# Use an official Go image to build the binary
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app, producing a binary called "app"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app server.go

# Use a minimal Docker image for the final stage
FROM scratch

# Copy the CA certificates from the builder stage (from /etc/ssl/certs in the Go image)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the binary from the builder stage
COPY --from=builder /app/app /app

EXPOSE 3000

# Command to run the executable
ENTRYPOINT ["/app"]
