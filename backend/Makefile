# Define project-specific variables
PROJECT_NAME := naytife
SQLC_DIR := internal/db
GQLGEN_CONFIG := gqlgen.yml
ATLAS_MIGRATIONS_DIR := internal/db/migrations
ATLAS_SCHEMA_DIR := internal/db/schema.sql
ATLAS_DEV_URL := docker://postgres?search_path=public

# Define Go-related variables
GO := go
GOFMT := gofmt
GOTEST := go test
GOMOD := go mod
GOBUILD := go build
GOLINT := golangci-lint run
AIR := air  --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"
# Define binary name
BINARY_NAME := bin/api

# Default target: run everything needed to start the project
.PHONY: all
all: generate build

# Code generation targets
.PHONY: generate
generate: generate-sqlc generate-gqlgen ## Run all code generators

.PHONY: generate-sqlc
generate-sqlc: ## Generate SQL code from .sql files
	sqlc generate

.PHONY: generate-gqlgen
generate-gqlgen: ## Generate GraphQL code from schema
	go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate 

# Migration targets
.PHONY: migrate-up
migrate-up: ## Apply all up migrations
	atlas migrate apply --env local

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	atlas migrate down --env local --dev-url $(ATLAS_DEV_URL)

.PHONY: migrate-clean
migrate-clean: ## Clean the database
	atlas schema clean --env local

.PHONY: migrate-new
migrate-new: ## Create a new migration file
	atlas migrate diff --dir file://$(ATLAS_MIGRATIONS_DIR) --to file://$(ATLAS_SCHEMA_DIR) --dev-url $(ATLAS_DEV_URL) $(name)

# Formatting and linting targets
.PHONY: fmt
fmt: ## Format the code
	$(GOFMT) -s -w .

.PHONY: lint
lint: ## Lint the code using golangci-lint
	$(GOLINT)

# Testing targets
.PHONY: test
test: ## Run all tests
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Build target
.PHONY: build
build: ## Build the Go binary
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/api

# Air target for auto-compilation
.PHONY: dev
dev: ## Run the project with Air for auto-compilation
	$(AIR)

# Clean target
.PHONY: clean
clean: ## Clean up generated files
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Help target
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
