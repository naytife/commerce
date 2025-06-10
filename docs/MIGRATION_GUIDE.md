# Database Migration Guide

This guide explains how to manage database migrations for the Naytife Commerce Platform.

## Overview

The migration system uses [Atlas](https://atlasgo.io/) for schema migrations with Kubernetes Jobs for deployment automation.

## Migration Files

Migration files are located in `backend/internal/db/migrations/` and follow the naming convention:
```
YYYYMMDDHHMMSS_description.sql
```

Example: `20250603120000_add_stock_movements.sql`

## Local Development

### Prerequisites
- Go installed
- PostgreSQL running locally
- Atlas CLI installed

### Common Commands

```bash
# Navigate to backend directory
cd backend

# Check migration status
make migrate-status

# Apply all pending migrations
make migrate-up

# Validate migration files
make migrate-validate

# Create a new migration
make migrate-new name=add_new_feature

# Rollback to specific version
make migrate-rollback version=20250603120000

# Reset database (clean + migrate)
make migrate-reset

# Dry run (test without applying)
make migrate-dry-run
```

## Kubernetes/Production

### Management Script

Use the migration management script for Kubernetes operations:

```bash
# Check current migration status
./k3s/scripts/manage-migrations.sh status

# Run migrations
./k3s/scripts/manage-migrations.sh run

# View migration logs
./k3s/scripts/manage-migrations.sh logs

# Rollback to specific version
./k3s/scripts/manage-migrations.sh rollback 20250603120000

# Validate migration files
./k3s/scripts/manage-migrations.sh validate

# Clean failed jobs
./k3s/scripts/manage-migrations.sh clean
```

### Makefile Integration

```bash
# Run migrations in Kubernetes
make k8s-migrate-run

# Check Kubernetes migration status
make k8s-migrate-status

# View Kubernetes migration logs
make k8s-migrate-logs

# Rollback in Kubernetes
make k8s-migrate-rollback version=20250603120000
```

## Migration Pipeline

For automated deployments, use the migration pipeline:

```bash
# Full migration pipeline with validation
./scripts/migration-pipeline.sh
```

This pipeline:
1. Validates migration files
2. Updates Kubernetes ConfigMap
3. Runs migration job
4. Monitors progress
5. Verifies success

## Writing Migrations

### Best Practices

1. **Always use transactions** (when possible)
2. **Make migrations reversible**
3. **Test migrations on copy of production data**
4. **Use descriptive names**
5. **Break large changes into smaller migrations**

### Example Migration

```sql
-- 20250608120000_add_user_preferences.sql
BEGIN;

-- Add user preferences table
CREATE TABLE user_preferences (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    preference_key VARCHAR(100) NOT NULL,
    preference_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id, preference_key)
);

-- Create index for faster lookups
CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

-- Add RLS (Row Level Security) for multi-tenancy
ALTER TABLE user_preferences ENABLE ROW LEVEL SECURITY;

-- RLS Policy (example)
CREATE POLICY user_preferences_policy ON user_preferences
    FOR ALL USING (user_id = current_setting('app.current_user_id', true));

COMMIT;
```

## Atlas Configuration

The Atlas configuration is in `backend/atlas.hcl`:

```hcl
# Local environment
env "local" {
  url = "postgres://naytife:postgres@localhost:5432/naytifedb?search_path=naytife_schema&sslmode=disable"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# Production environment
env "prod" {
  migration {
    dir = "file://internal/db/migrations"
  }
}
```

## Kubernetes Job Configuration

The unified migration job (`k3s/manifests/06-backend/backend-migration.yaml`) includes:

- **Wait for PostgreSQL**: Init container ensures database is ready
- **Pre-migration checks**: Validates database connectivity and schema
- **Migration execution**: Runs Atlas migrations with detailed logging
- **Resource limits**: Optimized CPU and memory constraints
- **Retry logic**: Up to 3 attempts on failure with backoff
- **Timeout**: 10-minute maximum execution time
- **Health monitoring**: Liveness probe to track migration progress
- **Utility scripts**: ConfigMap with helper scripts for manual operations

## Troubleshooting

### Common Issues

1. **Migration job stuck**
   ```bash
   # Check job status
   kubectl get jobs -n naytife
   
   # View logs
   kubectl logs job/backend-migrate -n naytife
   
   # Clean and restart
   ./k3s/scripts/manage-migrations.sh clean
   ./k3s/scripts/manage-migrations.sh run
   ```

2. **Migration validation failed**
   ```bash
   # Check file syntax
   ./k3s/scripts/manage-migrations.sh validate
   
   # Test locally first
   make migrate-validate
   ```

3. **Database connection issues**
   ```bash
   # Check PostgreSQL status
   kubectl get pods -n naytife -l app=postgres
   
   # Test connection
   kubectl run db-test --rm -i --restart=Never --image=postgres:15-alpine -n naytife -- \
     psql -h postgres.naytife.svc.cluster.local -U postgres -d naytifedb -c "SELECT version();"
   ```

### Recovery Procedures

#### Failed Migration Recovery

1. **Check the error**:
   ```bash
   kubectl logs job/backend-migrate -n naytife
   ```

2. **Fix the migration file** if needed

3. **Clean and retry**:
   ```bash
   ./k3s/scripts/manage-migrations.sh clean
   ./k3s/scripts/manage-migrations.sh run
   ```

#### Rollback Procedure

1. **Identify target version**:
   ```bash
   ./k3s/scripts/manage-migrations.sh status
   ```

2. **Perform rollback**:
   ```bash
   ./k3s/scripts/manage-migrations.sh rollback [version]
   ```

3. **Verify result**:
   ```bash
   ./k3s/scripts/manage-migrations.sh status
   ```

## Monitoring and Alerting

### Key Metrics to Monitor

- Migration job completion status
- Migration execution time
- Database connection health
- Migration file validation status

### Log Monitoring

Migration logs include:
- Pre-migration checks
- Migration execution progress
- Error details and stack traces
- Post-migration verification

### Health Checks

The migration job includes health checks:
- PostgreSQL connectivity
- Schema accessibility
- Migration file integrity

## Security Considerations

1. **Database credentials** are stored in Kubernetes secrets
2. **Migration files** are mounted read-only
3. **Row Level Security** is enabled for multi-tenancy
4. **Connection encryption** is enforced in production

## Performance Considerations

1. **Resource limits** prevent migration job from consuming excessive resources
2. **Timeout settings** prevent stuck migrations
3. **Index creation** is done with appropriate strategies
4. **Large data migrations** should be done in batches

## Integration with CI/CD

The migration system integrates with CI/CD pipelines:

1. **Pre-deployment validation**
2. **Automated migration execution**
3. **Post-deployment verification**
4. **Rollback capabilities**

For automated deployments, use:
```bash
./scripts/migration-pipeline.sh
```

This ensures consistent and reliable database deployments.
