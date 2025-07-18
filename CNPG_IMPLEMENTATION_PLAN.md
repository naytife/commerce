# CloudNativePG (CNPG) Implementation Plan

## Executive Summary

This plan outlines the greenfield migration from the current basic PostgreSQL Deployment to CloudNativePG (CNPG), a production-grade PostgreSQL operator that provides high availability, automated backups, monitoring, and operational excellence for your Naytife commerce platform.

## Current State Analysis

### Current PostgreSQL Setup
- **Deployment Type**: Basic Kubernetes Deployment with single replica
- **Image**: `postgres:15-alpine`
- **Storage**: EmptyDir (local), PVC (production) - 50Gi
- **Namespaces**: `naytife` (local), `naytife-staging`, `naytife-production`
- **Database**: `naytifedb` with schemas: `hydra`, `naytife_schema`
- **Connection**: Direct database connections via `DATABASE_URL`
- **Migrations**: Atlas-based schema migrations
- **Secrets**: SOPS-encrypted secrets with separate files per environment

### Current Limitations
- No high availability or automatic failover
- No automated backup/restore capabilities
- No connection pooling
- No monitoring/alerting
- No point-in-time recovery (PITR)
- Manual scaling and maintenance
- Basic resource management

## Target Architecture with CNPG

### Core Components
1. **CloudNativePG Operator**: Manages PostgreSQL clusters lifecycle
2. **PostgreSQL Cluster**: High-availability PostgreSQL cluster with replicas
3. **PgBouncer**: Connection pooling for improved performance
4. **Backup System**: Automated backups with PITR capabilities
5. **Monitoring**: Prometheus metrics and Grafana dashboards
6. **Storage**: Persistent volumes with proper storage classes

## Implementation Plan

### Phase 1: Prerequisites and Operator Installation

#### 1.1 Environment Preparation
```bash
# Prerequisites validation
- Kubernetes cluster with CSI storage support
- Helm 3.x installed
- kubectl configured
- Storage class available (local-path for local, oci for production)
- Resource requirements: 2+ CPU cores, 4GB+ RAM per environment
```

#### 1.2 CNPG Operator Installation
```yaml
# File: deploy/base/cnpg-operator/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cnpg-system
  labels:
    name: cnpg-system
```

```yaml
# File: deploy/base/cnpg-operator/helm-release.yaml
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: HelmRepository
metadata:
  name: cnpg
  namespace: cnpg-system
spec:
  interval: 5m
  url: https://cloudnative-pg.github.io/charts
---
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: cloudnative-pg
  namespace: cnpg-system
spec:
  interval: 5m
  chart:
    spec:
      chart: cloudnative-pg
      sourceRef:
        kind: HelmRepository
        name: cnpg
      interval: 5m
  values:
    crds:
      create: true
    monitoring:
      enabled: true
    config:
      data:
        INHERITED_ANNOTATIONS: "prometheus.io/scrape"
```

### Phase 2: Storage and Backup Configuration

#### 2.1 Storage Classes
```yaml
# File: deploy/base/cnpg-storage/storage-class.yaml
# For local development
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: cnpg-local-storage
provisioner: rancher.io/local-path
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Delete
---
# For production (OCI)
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: cnpg-oci-storage
provisioner: blockvolume.csi.oraclecloud.com
parameters:
  fsType: ext4
  replicaCount: "3"
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Retain
```

#### 2.2 Backup Storage Configuration
```yaml
# File: deploy/base/cnpg-backup/backup-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: cnpg-backup-secret
  namespace: naytife
type: Opaque
stringData:
  # For local: MinIO credentials
  ACCESS_KEY_ID: "minio-access-key"
  SECRET_ACCESS_KEY: "minio-secret-key"
  # For production: OCI Object Storage
  # ACCESS_KEY_ID: "oci-access-key"
  # SECRET_ACCESS_KEY: "oci-secret-key"
```

### Phase 3: Database Cluster Configuration

#### 3.1 Base PostgreSQL Cluster
```yaml
# File: deploy/base/cnpg-cluster/cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: naytife-postgres
  namespace: naytife
spec:
  instances: 3
  
  postgresql:
    parameters:
      # Performance tuning
      shared_preload_libraries: "pg_stat_statements"
      max_connections: "200"
      shared_buffers: "256MB"
      effective_cache_size: "1GB"
      maintenance_work_mem: "64MB"
      checkpoint_completion_target: "0.9"
      wal_buffers: "16MB"
      default_statistics_target: "100"
      random_page_cost: "1.1"
      effective_io_concurrency: "200"
      work_mem: "4MB"
      
      # Logging
      log_statement: "all"
      log_min_duration_statement: "1000"
      log_checkpoints: "on"
      log_connections: "on"
      log_disconnections: "on"
      log_lock_waits: "on"

  bootstrap:
    initdb:
      database: naytifedb
      owner: naytife
      secret:
        name: cnpg-cluster-secret
      postInitTemplateSQL:
        - CREATE SCHEMA IF NOT EXISTS hydra;
        - CREATE SCHEMA IF NOT EXISTS naytife_schema;
        - GRANT ALL PRIVILEGES ON SCHEMA hydra TO naytife;
        - GRANT ALL PRIVILEGES ON SCHEMA naytife_schema TO naytife;

  storage:
    size: 50Gi
    storageClass: cnpg-local-storage
    
  monitoring:
    enabled: true
    
  nodeMaintenanceWindow:
    inProgress: false
    reusePVC: true
```

#### 3.2 Environment-Specific Overlays

**Local Environment:**
```yaml
# File: deploy/overlays/local/cnpg-cluster-local-patch.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: naytife-postgres
  namespace: naytife
spec:
  instances: 1  # Single instance for local development
  
  storage:
    size: 10Gi
    storageClass: cnpg-local-storage
    
  resources:
    requests:
      memory: "256Mi"
      cpu: "200m"
    limits:
      memory: "512Mi"
      cpu: "500m"
```

**Production Environment:**
```yaml
# File: deploy/overlays/production/cnpg-cluster-production-patch.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: naytife-postgres
  namespace: naytife-production
spec:
  instances: 3  # High availability
  
  storage:
    size: 100Gi
    storageClass: cnpg-oci-storage
    
  resources:
    requests:
      memory: "1Gi"
      cpu: "500m"
    limits:
      memory: "2Gi"
      cpu: "1000m"
      
  backup:
    barmanObjectStore:
      destinationPath: "s3://naytife-postgres-backups/production"
      endpointURL: "https://objectstorage.us-phoenix-1.oraclecloud.com"
      data:
        retention: "7d"
        compression: "gzip"
      wal:
        retention: "1d"
        compression: "gzip"
      s3Credentials:
        accessKeyId:
          name: cnpg-backup-secret
          key: ACCESS_KEY_ID
        secretAccessKey:
          name: cnpg-backup-secret
          key: SECRET_ACCESS_KEY
```

### Phase 4: Connection Pooling with PgBouncer

#### 4.1 PgBouncer Configuration
```yaml
# File: deploy/base/cnpg-pooler/pooler.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Pooler
metadata:
  name: naytife-postgres-pooler
  namespace: naytife
spec:
  cluster:
    name: naytife-postgres
  
  instances: 2
  type: rw
  
  pgbouncer:
    poolMode: transaction
    parameters:
      max_client_conn: "100"
      default_pool_size: "20"
      max_db_connections: "25"
      max_user_connections: "25"
      server_reset_query: "DISCARD ALL"
      server_check_query: "SELECT 1"
      server_check_delay: "10"
      
  monitoring:
    enabled: true
    
  resources:
    requests:
      memory: "64Mi"
      cpu: "50m"
    limits:
      memory: "128Mi"
      cpu: "100m"
```

### Phase 5: Monitoring and Alerting

#### 5.1 ServiceMonitor for Prometheus
```yaml
# File: deploy/base/cnpg-monitoring/service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnpg-cluster-metrics
  namespace: naytife
spec:
  selector:
    matchLabels:
      postgresql: naytife-postgres
      role: primary
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cnpg-pooler-metrics
  namespace: naytife
spec:
  selector:
    matchLabels:
      cnpg.io/poolerName: naytife-postgres-pooler
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

#### 5.2 Alerting Rules
```yaml
# File: deploy/base/cnpg-monitoring/alerts.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: cnpg-alerts
  namespace: naytife
spec:
  groups:
  - name: cnpg.rules
    rules:
    - alert: PostgreSQLInstanceDown
      expr: pg_up == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "PostgreSQL instance is down"
        description: "PostgreSQL instance {{ $labels.instance }} is down"
        
    - alert: PostgreSQLReplicationLag
      expr: pg_stat_replication_lag > 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "PostgreSQL replication lag is high"
        description: "Replication lag is {{ $value }}s on {{ $labels.instance }}"
        
    - alert: PostgreSQLConnectionsHigh
      expr: pg_stat_database_numbackends / pg_settings_max_connections > 0.8
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "PostgreSQL connections are high"
        description: "Connection usage is {{ $value }}% on {{ $labels.instance }}"
```

### Phase 6: Migration Strategy

#### 6.1 Data Migration Process
```bash
# Migration script: deploy/scripts/cnpg-migration.sh
#!/bin/bash

set -e

ENVIRONMENT=${1:-local}
NAMESPACE="naytife"
[ "$ENVIRONMENT" != "local" ] && NAMESPACE="naytife-$ENVIRONMENT"

print_header() {
    echo -e "\033[0;36m📋 $1\033[0m"
}

print_info() {
    echo -e "\033[0;34mℹ️ $1\033[0m"
}

print_success() {
    echo -e "\033[0;32m✅ $1\033[0m"
}

print_error() {
    echo -e "\033[0;31m❌ $1\033[0m"
}

# Pre-migration validation
validate_environment() {
    print_header "Pre-migration validation"
    
    # Check if CNPG operator is installed
    if ! kubectl get crd clusters.postgresql.cnpg.io >/dev/null 2>&1; then
        print_error "CNPG operator not installed"
        exit 1
    fi
    
    # Check current postgres pod
    if ! kubectl get deployment postgres -n "$NAMESPACE" >/dev/null 2>&1; then
        print_error "Current postgres deployment not found"
        exit 1
    fi
    
    print_success "Environment validation passed"
}

# Create database dump
create_backup() {
    print_header "Creating database backup"
    
    # Get current postgres pod
    POSTGRES_POD=$(kubectl get pods -n "$NAMESPACE" -l app=postgres -o jsonpath='{.items[0].metadata.name}')
    
    if [ -z "$POSTGRES_POD" ]; then
        print_error "No postgres pod found"
        exit 1
    fi
    
    # Create backup
    kubectl exec -n "$NAMESPACE" "$POSTGRES_POD" -- pg_dump -U naytife naytifedb > "/tmp/naytife-backup-$(date +%Y%m%d_%H%M%S).sql"
    
    print_success "Backup created successfully"
}

# Deploy CNPG cluster
deploy_cnpg() {
    print_header "Deploying CNPG cluster"
    
    # Apply CNPG manifests
    kubectl apply -k "deploy/overlays/$ENVIRONMENT/cnpg"
    
    # Wait for cluster to be ready
    kubectl wait --for=condition=Ready cluster/naytife-postgres -n "$NAMESPACE" --timeout=300s
    
    print_success "CNPG cluster deployed successfully"
}

# Migrate data
migrate_data() {
    print_header "Migrating data to CNPG"
    
    # Get CNPG primary pod
    CNPG_PRIMARY=$(kubectl get pods -n "$NAMESPACE" -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -z "$CNPG_PRIMARY" ]; then
        print_error "CNPG primary pod not found"
        exit 1
    fi
    
    # Restore backup
    kubectl exec -n "$NAMESPACE" "$CNPG_PRIMARY" -- psql -U naytife -d naytifedb < "/tmp/naytife-backup-$(date +%Y%m%d_%H%M%S).sql"
    
    print_success "Data migration completed"
}

# Update application configuration
update_app_config() {
    print_header "Updating application configuration"
    
    # Update DATABASE_URL to use pooler
    NEW_DATABASE_URL="postgresql://naytife:${POSTGRES_PASSWORD}@naytife-postgres-pooler-rw.${NAMESPACE}.svc.cluster.local:5432/naytifedb?sslmode=require"
    
    print_info "Update DATABASE_URL in secrets to: $NEW_DATABASE_URL"
    print_info "Manual step: Update deploy/secrets/${ENVIRONMENT}/backend-secret.yaml"
    
    print_success "Configuration update guidance provided"
}

# Validation and cleanup
validate_migration() {
    print_header "Validating migration"
    
    # Test database connectivity
    CNPG_PRIMARY=$(kubectl get pods -n "$NAMESPACE" -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    # Test connection through pooler
    kubectl exec -n "$NAMESPACE" "$CNPG_PRIMARY" -- psql -U naytife -d naytifedb -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema IN ('hydra', 'naytife_schema');"
    
    print_success "Migration validation passed"
}

# Main execution
main() {
    print_header "CNPG Migration for $ENVIRONMENT environment"
    
    validate_environment
    create_backup
    deploy_cnpg
    migrate_data
    update_app_config
    validate_migration
    
    print_success "Migration completed successfully!"
    print_info "Next steps:"
    print_info "  1. Update application secrets with new connection string"
    print_info "  2. Restart application deployments"
    print_info "  3. Verify application connectivity"
    print_info "  4. Remove old postgres deployment when stable"
}

main "$@"
```

### Phase 7: Application Integration

#### 7.1 Updated Secret Configuration
```yaml
# File: deploy/secrets/local/cnpg-cluster-secret.yaml (encrypted with SOPS)
apiVersion: v1
kind: Secret
metadata:
  name: cnpg-cluster-secret
  namespace: naytife
type: postgresql.cnpg.io/ClusterSecret
stringData:
  username: naytife
  password: your-secure-password-here
  # Auto-generated by CNPG:
  # postgres: superuser-password
  # app: app-user-password
```

#### 7.2 Backend Configuration Update
```yaml
# File: deploy/secrets/local/backend-secret.yaml (updated)
apiVersion: v1
kind: Secret
metadata:
  name: backend-secret
  namespace: naytife
type: Opaque
stringData:
  # Updated to use pooler service
  DATABASE_URL: "postgresql://naytife:${POSTGRES_PASSWORD}@naytife-postgres-pooler-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
  # Direct connection (for migrations)
  DATABASE_URL_DIRECT: "postgresql://naytife:${POSTGRES_PASSWORD}@naytife-postgres-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
  # Other secrets remain the same
```

### Phase 8: Backup and Recovery

#### 8.1 Backup Schedule
```yaml
# File: deploy/base/cnpg-backup/scheduled-backup.yaml
apiVersion: postgresql.cnpg.io/v1
kind: ScheduledBackup
metadata:
  name: naytife-postgres-backup
  namespace: naytife
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  backupOwnerReference: self
  cluster:
    name: naytife-postgres
```

#### 8.2 Recovery Procedures
```bash
# Recovery script: deploy/scripts/cnpg-recovery.sh
#!/bin/bash

# Point-in-time recovery
restore_to_point_in_time() {
    local target_time=$1
    
    cat <<EOF | kubectl apply -f -
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: naytife-postgres-recovered
  namespace: naytife
spec:
  instances: 1
  
  bootstrap:
    recovery:
      backup:
        name: naytife-postgres-backup
      recoveryTarget:
        targetTime: "$target_time"
      
  storage:
    size: 50Gi
    storageClass: cnpg-local-storage
EOF
}

# Full backup restore
restore_from_backup() {
    local backup_name=$1
    
    cat <<EOF | kubectl apply -f -
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: naytife-postgres-restored
  namespace: naytife
spec:
  instances: 1
  
  bootstrap:
    recovery:
      backup:
        name: "$backup_name"
      
  storage:
    size: 50Gi
    storageClass: cnpg-local-storage
EOF
}
```

### Phase 9: Testing and Validation

#### 9.1 Integration Tests Update
```bash
# File: deploy/scripts/test-cnpg-integration.sh
#!/bin/bash

test_cnpg_cluster() {
    local namespace=$1
    local cluster_name="naytife-postgres"
    
    print_header "Testing CNPG Cluster"
    
    # Test cluster status
    run_test "CNPG cluster is ready" "kubectl get cluster $cluster_name -n $namespace -o jsonpath='{.status.phase}' | grep -q 'Cluster in healthy state'"
    
    # Test primary pod
    run_test "Primary pod is running" "kubectl get pods -n $namespace -l postgresql=$cluster_name,role=primary -o jsonpath='{.items[0].status.phase}' | grep -q 'Running'"
    
    # Test pooler
    run_test "Pooler is ready" "kubectl get pooler naytife-postgres-pooler -n $namespace -o jsonpath='{.status.phase}' | grep -q 'Ready'"
    
    # Test database connectivity
    PRIMARY_POD=$(kubectl get pods -n $namespace -l postgresql=$cluster_name,role=primary -o jsonpath='{.items[0].metadata.name}')
    run_test "Database connection" "kubectl exec -n $namespace $PRIMARY_POD -- psql -U naytife -d naytifedb -c 'SELECT 1;'"
    
    # Test pooler connectivity
    run_test "Pooler connection" "kubectl exec -n $namespace $PRIMARY_POD -- psql -h naytife-postgres-pooler-rw -U naytife -d naytifedb -c 'SELECT 1;'"
    
    print_success "CNPG cluster tests passed"
}
```

### Phase 10: Monitoring and Observability

#### 10.1 Grafana Dashboard
```json
{
  "dashboard": {
    "title": "CloudNativePG - Naytife PostgreSQL",
    "panels": [
      {
        "title": "Cluster Status",
        "type": "stat",
        "targets": [
          {
            "expr": "cnpg_cluster_status{cluster=\"naytife-postgres\"}",
            "legendFormat": "Status"
          }
        ]
      },
      {
        "title": "Connection Count",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_database_numbackends{cluster=\"naytife-postgres\"}",
            "legendFormat": "Connections"
          }
        ]
      },
      {
        "title": "Replication Lag",
        "type": "graph",
        "targets": [
          {
            "expr": "pg_stat_replication_lag{cluster=\"naytife-postgres\"}",
            "legendFormat": "Lag (seconds)"
          }
        ]
      }
    ]
  }
}
```

## Implementation Timeline

### Week 1: Foundation
- Install CNPG operator
- Set up storage classes
- Configure basic cluster for local environment
- Test basic functionality

### Week 2: Production Setup
- Configure production storage and backup
- Set up monitoring and alerting
- Create migration scripts
- Test disaster recovery procedures

### Week 3: Migration
- Perform local environment migration
- Test application integration
- Validate performance and functionality
- Create staging environment

### Week 4: Production Deployment
- Migrate staging environment
- Comprehensive testing
- Production migration (with rollback plan)
- Documentation and training

## Risk Mitigation

### Technical Risks
1. **Data Loss**: Comprehensive backup strategy with point-in-time recovery
2. **Downtime**: Blue-green deployment approach with connection string updates
3. **Performance**: Extensive load testing and connection pooling
4. **Compatibility**: Thorough testing of all application features

### Operational Risks
1. **Team Training**: Comprehensive documentation and training sessions
2. **Monitoring**: Detailed alerting and monitoring setup
3. **Rollback Plan**: Maintain old deployment until stable
4. **Support**: Establish incident response procedures

## Success Metrics

### Technical Metrics
- **RTO (Recovery Time Objective)**: < 15 minutes
- **RPO (Recovery Point Objective)**: < 1 minute
- **Availability**: > 99.9%
- **Performance**: Query response time < 100ms (95th percentile)

### Operational Metrics
- **Automated Backups**: 100% success rate
- **Monitoring Coverage**: 100% of critical metrics
- **Incident Response**: < 5 minutes mean time to detection
- **Documentation**: 100% of procedures documented

## Post-Implementation

### Ongoing Maintenance
1. **Regular Updates**: Monthly CNPG operator updates
2. **Backup Validation**: Weekly restore testing
3. **Performance Tuning**: Monthly performance reviews
4. **Security Updates**: Immediate security patches

### Future Enhancements
1. **Multi-region Setup**: Cross-region replication
2. **Advanced Monitoring**: Custom metrics and dashboards
3. **Automated Scaling**: HPA for connection pools
4. **Disaster Recovery**: Multi-cloud backup strategy

## Directory Structure Changes

The implementation will add the following directories to your existing structure:

```
deploy/
├── base/
│   ├── cnpg-operator/
│   │   ├── namespace.yaml
│   │   └── helm-release.yaml
│   ├── cnpg-storage/
│   │   └── storage-class.yaml
│   ├── cnpg-cluster/
│   │   └── cluster.yaml
│   ├── cnpg-pooler/
│   │   └── pooler.yaml
│   ├── cnpg-backup/
│   │   ├── backup-secret.yaml
│   │   └── scheduled-backup.yaml
│   └── cnpg-monitoring/
│       ├── service-monitor.yaml
│       └── alerts.yaml
├── overlays/
│   ├── local/
│   │   └── cnpg-cluster-local-patch.yaml
│   ├── staging/
│   │   └── cnpg-cluster-staging-patch.yaml
│   └── production/
│       └── cnpg-cluster-production-patch.yaml
├── scripts/
│   ├── cnpg-migration.sh
│   ├── cnpg-recovery.sh
│   └── test-cnpg-integration.sh
└── secrets/
    ├── local/
    │   └── cnpg-cluster-secret.yaml
    ├── staging/
    │   └── cnpg-cluster-secret.yaml
    └── production/
        └── cnpg-cluster-secret.yaml
```

This comprehensive plan provides a production-ready PostgreSQL solution with high availability, automated backups, monitoring, and operational excellence while maintaining compatibility with your existing application architecture.
