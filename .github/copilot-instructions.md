# Naytife Commerce Platform - AI Assistant Instructions

## Architecture Overview

This is a **multi-tenant e-commerce platform** built with Go backend, SvelteKit frontends, and Kubernetes deployment. The system consists of 8+ microservices orchestrated via k3s/Kubernetes with GitOps deployment patterns.

### Core Technology Stack
- **Backend**: Go with Fiber framework + GraphQL (gqlgen)
- **Database**: PostgreSQL with SQLC for type-safe queries + Atlas for migrations
- **Frontend**: SvelteKit (3 separate apps: dashboard, website, templates)
- **Auth**: Ory Stack (Hydra OAuth2, Oathkeeper gateway, Keto access control)
- **Deployment**: Kustomize + SOPS + Skaffold for k3s/Kubernetes
- **Payments**: Multi-provider (Stripe, PayPal, Paystack, Flutterwave)

## Critical Development Workflows

### Backend Development
```bash
# Code generation (run after schema/query changes)
make generate  # Runs both SQLC and GraphQL generation
make generate-sqlc  # Generate type-safe DB code from internal/db/queries/*.sql
make generate-gqlgen  # Generate GraphQL from internal/gql/public/schema/*.graphql

# Development server with hot reload
make dev  # Uses Air for hot reloading

# Database migrations
atlas migrate apply --env local  # Apply migrations locally
atlas migrate diff --env local   # Generate new migration
```

### Deployment Patterns
```bash
# Complete platform deployment (preferred for development)
./deploy/scripts/deploy.sh local

# Individual service development
cd backend && make dev
cd dashboard && npm run dev
cd templates/template_1 && npm run dev

# Skaffold development (live reload in k3s)
skaffold dev --port-forward
```

## Project-Specific Conventions

### Database Patterns
- **Multi-tenancy**: All tables include `shop_id` for tenant isolation
- **Type Safety**: Use SQLC - write SQL in `internal/db/queries/`, generate Go code
- **Migrations**: Atlas manages schema - edit `internal/db/schema.sql`, run `atlas migrate diff`
- **Connection**: Uses pgx/v5 with connection pooling

### GraphQL Architecture
- **Schema**: Located in `internal/gql/public/schema/*.graphql`
- **Resolvers**: Follow-schema layout in `internal/gql/public/resolver/`
- **Models**: Auto-generated in `internal/gql/public/model/models_gen.go`
- **Config**: gqlgen.yml defines generation rules

### Authentication Flow
- **OAuth2**: Ory Hydra handles all OAuth flows
- **Gateway**: Oathkeeper validates requests and injects user context
- **Handler**: `auth/authentication-handler/` bridges social logins to Hydra
- **Frontend**: Uses `$page.data.user` for auth state in SvelteKit

### Kubernetes Deployment
- **Structure**: `deploy/base/` contains Kustomize base configs
- **Environments**: `deploy/overlays/{local,staging,production}/` for environment-specific overrides
- **Secrets**: SOPS-encrypted in `deploy/secrets/` - use `deploy/scripts/sops-helper.sh`
- **Services**: Each microservice has its own base config and overlay

## Key Integration Points

### Service Communication
- **API Gateway**: Oathkeeper routes and authenticates all requests
- **Backend**: Exposes REST (`/api/v1/`) and GraphQL (`/graphql`) endpoints
- **Frontend Auth**: SvelteKit `hooks.server.ts` handles auth flow with Hydra
- **Database**: Single PostgreSQL instance with RLS for multi-tenancy

### Payment Processing
- **Multi-provider**: Backend abstracts Stripe, PayPal, Paystack, Flutterwave
- **Shop Config**: Payment methods stored in `shop_payment_methods` table
- **Processing**: Order creation triggers payment flow based on shop configuration

### Template System
- **Registry**: `services/template-registry/` manages storefront templates
- **Deployer**: `services/store-deployer/` provisions new tenant storefronts
- **Templates**: Located in `templates/` - template_1 is default SvelteKit storefront

## Development Guidelines

### Adding New Features
1. **Backend**: Add SQL queries to `internal/db/queries/`, run `make generate-sqlc`
2. **GraphQL**: Update schema in `internal/gql/public/schema/`, run `make generate-gqlgen`
3. **Frontend**: Use TypeScript, follow existing patterns in `src/routes/`
4. **Deployment**: Add Kustomize configs to `deploy/base/` and environment overlays

### Testing Locally
- **Full Stack**: Use `./deploy/scripts/deploy.sh local` for complete environment
- **Backend Only**: `make dev` runs with Air hot reload
- **Frontend Only**: `npm run dev` with TLS disabled (`NODE_TLS_REJECT_UNAUTHORIZED=0`)

### Common Gotchas
- **CSRF**: Dashboard has `csrf: { checkOrigin: false }` for Ory integration
- **TLS**: Local development uses self-signed certs - set `NODE_TLS_REJECT_UNAUTHORIZED=0`
- **Migrations**: Always test migrations with `atlas migrate apply --dry-run`
- **Secrets**: Use `deploy/scripts/sops-helper.sh` to manage encrypted secrets

### File Patterns to Recognize
- `*.sql.go`: SQLC-generated database code (don't edit manually)
- `*_gen.go`: GraphQL-generated code (don't edit manually)
- `kustomization.yaml`: Kustomize configuration files
- `*.sops.yaml`: SOPS-encrypted secrets
