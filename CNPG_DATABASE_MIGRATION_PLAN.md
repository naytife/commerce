# CNPG Database Connection Migration Plan

## Overview

This plan outlines the migration strategy for updating database connections to use CloudNativePG (CNPG) pooled connections. The approach is simple: **no backend code changes needed** - services just use whatever `DATABASE_URL` is provided in their environment.

## Key Principles

1. **No backend Go code changes needed** - applications use the provided `DATABASE_URL`
2. **Migration jobs use direct connections** - for schema operations
3. **Runtime operations use pooled connections** - for better performance
4. **Atlas configuration is for local development** - when developers run `make migrate-up`

## Connection Strategy

### Runtime Operations (Pooled Connections)
- **Service**: `naytife-postgres-pooler-rw.{namespace}.svc.cluster.local:5432`
- **Used by**: Backend service, Hydra service (runtime operations)
- **Benefits**: Connection pooling, load balancing, better performance

### Migration Operations (Direct Connections)
- **Service**: `naytife-postgres-rw.{namespace}.svc.cluster.local:5432`
- **Used by**: Migration jobs, Atlas CLI (local development)
- **Benefits**: Direct access for DDL operations, administrative tasks

## Current Database Consumers

| Service | Current Connection | New Connection | Update Required |
|---------|-------------------|----------------|-----------------|
| Backend Service | `postgres.naytife.svc.cluster.local:5432` | `naytife-postgres-pooler-rw` | **Secret only** |
| Hydra Service | `postgres.naytife.svc.cluster.local:5432` | `naytife-postgres-pooler-rw` | **Secret only** |
| Backend Migration Job | `postgres.naytife.svc.cluster.local:5432` | `naytife-postgres-rw` | **Job config only** |
| Hydra Migration Job | `postgres.naytife.svc.cluster.local:5432` | `naytife-postgres-rw` | **Job config only** |

## Implementation Plan

### Phase 1: Update Secret Configurations (Runtime Services)

#### 1.1 Backend Service Secret

**Files to Update:**
- `deploy/secrets/local/backend-secret.yaml`
- `deploy/secrets/staging/backend-secret.yaml`
- `deploy/secrets/production/backend-secret.yaml`

**Change:**
```yaml
stringData:
  # Update to use pooled connection
  DATABASE_URL: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
  # Keep other secrets unchanged
```

#### 1.2 Hydra Service Secret

**Files to Update:**
- `deploy/secrets/local/hydra-secret.yaml`
- `deploy/secrets/staging/hydra-secret.yaml`
- `deploy/secrets/production/hydra-secret.yaml`

**Change:**
```yaml
stringData:
  # Update to use pooled connection
  dsn: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=hydra"
```

### Phase 2: Update Migration Jobs (Direct Connections)

#### 2.1 Backend Migration Job

**File:** `deploy/overlays/local/backend-migrate-local-patch.yaml`

**Changes:**
```yaml
spec:
  template:
    spec:
      initContainers:
      - name: wait-for-postgres
        command:
        - sh
        - -c
        - |
          echo "Waiting for CNPG cluster to be ready..."
          until pg_isready -h naytife-postgres-rw.naytife.svc.cluster.local -p 5432 -U naytife; do
            echo "CNPG cluster not ready, waiting 2 seconds..."
            sleep 2
          done
          echo "✅ CNPG cluster is ready!"
      
      - name: migration-pre-check
        command:
        - sh
        - -c
        - |
          echo "🔍 Running pre-migration checks..."
          
          # Check database accessibility using direct connection
          PGPASSWORD=$POSTGRES_PASSWORD psql -h naytife-postgres-rw.naytife.svc.cluster.local -U naytife -d naytifedb -c "SELECT version();" || exit 1
          
          # Ensure schemas exist
          PGPASSWORD=$POSTGRES_PASSWORD psql -h naytife-postgres-rw.naytife.svc.cluster.local -U naytife -d naytifedb -c "CREATE SCHEMA IF NOT EXISTS naytife_schema;" || exit 1
          PGPASSWORD=$POSTGRES_PASSWORD psql -h naytife-postgres-rw.naytife.svc.cluster.local -U naytife -d naytifedb -c "CREATE SCHEMA IF NOT EXISTS hydra;" || exit 1
          
          echo "✅ Pre-migration checks passed"
      
      containers:
      - name: backend-migrate
        env:
        - name: DATABASE_URL
          value: "postgresql://naytife:${PASSWORD}@naytife-postgres-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
```

#### 2.2 Hydra Migration Job

**File:** `deploy/overlays/local/hydra-migrate-local-patch.yaml`

**Changes:**
```yaml
spec:
  template:
    spec:
      initContainers:
      - name: wait-for-postgres
        command:
        - sh
        - -c
        - |
          until pg_isready -h naytife-postgres-rw.naytife.svc.cluster.local -p 5432 -U naytife; do
            echo "Waiting for CNPG cluster..."
            sleep 2
          done
      
      containers:
      - name: hydra-migrate
        env:
        - name: DSN
          value: "postgresql://naytife:${PASSWORD}@naytife-postgres-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=hydra"
```

### Phase 3: Update Deployment Init Containers

#### 3.1 Backend Deployment

**File:** `deploy/overlays/local/backend-local-patch.yaml`

**Changes:**
```yaml
initContainers:
- name: wait-for-postgres
  image: postgres:15-alpine
  command:
  - sh
  - -c
  - |
    until pg_isready -h naytife-postgres-pooler-rw.naytife.svc.cluster.local -p 5432 -U naytife; do
      echo "Waiting for CNPG pooler..."
      sleep 2
    done
```

### Phase 4: Update Atlas Configuration (Local Development)

#### 4.1 Atlas Configuration

**File:** `backend/atlas.hcl`

**Changes:**
```hcl
# The "local" environment (traditional PostgreSQL for backward compatibility)
env "local" {
  url = "postgres://naytife:postgres@localhost:5432/naytifedb?search_path=naytife_schema&sslmode=disable"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# The "local_cnpg" environment (CNPG direct connection for local development)
env "local_cnpg" {
  url = "postgres://naytife:postgres@naytife-postgres-rw.naytife.svc.cluster.local:5432/naytifedb?search_path=naytife_schema&sslmode=require"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# The "prod" environment
env "prod" {
  migration {
    dir = "file://internal/db/migrations"
  }
}
```

#### 4.2 Makefile Update

**File:** `backend/Makefile`

**Add new target:**
```makefile
# CNPG Migration targets
.PHONY: migrate-up-cnpg
migrate-up-cnpg: ## Apply migrations to CNPG cluster (local development)
	atlas migrate apply --env local_cnpg

.PHONY: migrate-status-cnpg
migrate-status-cnpg: ## Show migration status on CNPG cluster
	atlas migrate status --env local_cnpg
```

## Environment-Specific Connection Strings

### Local Environment

**Backend Secret:**
```yaml
stringData:
  DATABASE_URL: "postgresql://naytife:postgres@naytife-postgres-pooler-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
```

**Hydra Secret:**
```yaml
stringData:
  dsn: "postgresql://naytife:postgres@naytife-postgres-pooler-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=hydra"
```

### Staging Environment

**Backend Secret:**
```yaml
stringData:
  DATABASE_URL: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife-staging.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
```

**Hydra Secret:**
```yaml
stringData:
  dsn: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife-staging.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=hydra"
```

### Production Environment

**Backend Secret:**
```yaml
stringData:
  DATABASE_URL: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife-production.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
```

**Hydra Secret:**
```yaml
stringData:
  dsn: "postgresql://naytife:${PASSWORD}@naytife-postgres-pooler-rw.naytife-production.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=hydra"
```

## What Does NOT Change

### Backend Go Code
- **No changes needed** - applications use whatever `DATABASE_URL` is provided
- **No new environment variables** - existing `DATABASE_URL` and `dsn` are sufficient
- **No repository modifications** - existing database layer works as-is

### Application Logic
- **No query changes** - all existing queries work the same
- **No connection handling changes** - existing connection pooling in Go continues to work
- **No schema changes** - database schemas remain unchanged

## Testing Strategy

### Phase 1: Local Environment

1. **Deploy CNPG cluster**
   ```bash
   kubectl apply -k deploy/base/cnpg-operator
   kubectl apply -k deploy/overlays/local
   ```

2. **Update secrets only**
   ```bash
   # Update backend-secret.yaml and hydra-secret.yaml
   # Re-encrypt with SOPS
   ```

3. **Test services restart with new connections**
   ```bash
   kubectl rollout restart deployment/backend -n naytife
   kubectl rollout restart deployment/hydra -n naytife-auth
   ```

4. **Verify connectivity**
   ```bash
   kubectl logs deployment/backend -n naytife
   kubectl logs deployment/hydra -n naytife-auth
   ```

### Phase 2: Migration Jobs

1. **Test migration jobs**
   ```bash
   kubectl delete job backend-migrate -n naytife
   kubectl delete job hydra-migrate -n naytife-auth
   kubectl apply -k deploy/overlays/local
   ```

2. **Verify migration success**
   ```bash
   kubectl logs job/backend-migrate -n naytife
   kubectl logs job/hydra-migrate -n naytife-auth
   ```

### Phase 3: Application Testing

1. **Test API endpoints**
   ```bash
   kubectl port-forward svc/backend 8000:8000 -n naytife
   curl http://localhost:8000/health
   ```

2. **Test OAuth2 flow**
   ```bash
   kubectl port-forward svc/hydra-public 4444:4444 -n naytife-auth
   curl http://localhost:4444/health/ready
   ```

## Rollback Strategy

### Quick Rollback

1. **Revert secret configurations**
   ```bash
   # Change DATABASE_URL back to: postgres.naytife.svc.cluster.local:5432
   # Change dsn back to: postgres.naytife.svc.cluster.local:5432
   ```

2. **Restart services**
   ```bash
   kubectl rollout restart deployment/backend -n naytife
   kubectl rollout restart deployment/hydra -n naytife-auth
   ```

3. **Verify rollback**
   ```bash
   kubectl logs deployment/backend -n naytife
   kubectl logs deployment/hydra -n naytife-auth
   ```

## Benefits

### Performance Benefits
- **Connection pooling** reduces connection overhead
- **Load balancing** improves concurrent request handling
- **Better resource utilization** with PgBouncer

### Operational Benefits
- **High availability** with automatic failover
- **Centralized connection management** through CNPG
- **Better monitoring** with connection pool metrics

### Maintenance Benefits
- **Simplified scaling** with connection pool management
- **Reduced connection limits** pressure on PostgreSQL
- **Improved backup and recovery** with CNPG

## Timeline

- **Week 1**: Update local environment secrets and test
- **Week 2**: Update migration jobs and test schema operations
- **Week 3**: Deploy to staging and validate
- **Week 4**: Deploy to production with monitoring

## Summary

This migration is straightforward:

1. **Services**: Change `DATABASE_URL` and `dsn` to use pooled connections
2. **Migration Jobs**: Change connection strings to use direct connections
3. **Atlas Config**: Add CNPG environment for local development
4. **No Go Code Changes**: Applications work with any valid connection string

The key insight is that the application layer is connection-agnostic - it just uses whatever database URL is provided in the environment.
