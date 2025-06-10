# Naytife Backend API

Go-based REST API with GraphQL support for the Naytife Commerce Platform.

## Features

- 🔗 REST API with Fiber framework
- 📊 GraphQL API with gqlgen
- 🗄️ PostgreSQL database with SQLC
- 🔐 OAuth2 authentication integration
- 📖 Swagger/OpenAPI documentation
- 🧪 Comprehensive test coverage

## Development

```bash
# Run locally
go run cmd/api/main.go

# Generate mocks
mockgen -source=./internal/db/repository.go -destination=./internal/mocks/repository_mock.go -package=mocks

# Run tests
go test ./...

# Build binary
make build
```

## API Documentation

When running locally, API docs are available at:
- Swagger UI: http://localhost:8000/docs
- GraphQL Playground: http://localhost:8000/graphql

## Deployment

For complete deployment instructions, see the [main deployment guide](../DEPLOYMENT_GUIDE.md).
