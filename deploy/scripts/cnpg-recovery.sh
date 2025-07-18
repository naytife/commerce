#!/bin/bash

# CNPG Recovery Script for Naytife Commerce Platform
# This script handles point-in-time recovery and backup restoration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

print_header() {
    echo -e "${CYAN}🔄 $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

show_usage() {
    echo -e "${CYAN}🔄 CNPG Recovery Script${NC}"
    echo "========================="
    echo "Usage: $0 [environment] [recovery-type] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo -e "${YELLOW}Recovery Types:${NC}"
    echo "  pitr              Point-in-time recovery"
    echo "  backup            Restore from backup"
    echo "  clone             Clone from existing cluster"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --target-time=TIME     Target time for PITR (ISO format)"
    echo "  --backup-name=NAME     Backup name for restore"
    echo "  --source-cluster=NAME  Source cluster for cloning"
    echo "  --new-cluster=NAME     New cluster name (default: naytife-postgres-recovered)"
    echo "  --dry-run             Show what would be done without executing"
    echo "  --help, -h            Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local pitr --target-time=2024-07-15T10:30:00Z"
    echo "  $0 staging backup --backup-name=backup-20240715-123456"
    echo "  $0 production clone --source-cluster=naytife-postgres"
}

# Default values
ENVIRONMENT=""
RECOVERY_TYPE=""
TARGET_TIME=""
BACKUP_NAME=""
SOURCE_CLUSTER=""
NEW_CLUSTER="naytife-postgres-recovered"
DRY_RUN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        pitr|backup|clone)
            RECOVERY_TYPE="$1"
            shift
            ;;
        --target-time=*)
            TARGET_TIME="${1#*=}"
            shift
            ;;
        --backup-name=*)
            BACKUP_NAME="${1#*=}"
            shift
            ;;
        --source-cluster=*)
            SOURCE_CLUSTER="${1#*=}"
            shift
            ;;
        --new-cluster=*)
            NEW_CLUSTER="${1#*=}"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check required arguments
if [ -z "$ENVIRONMENT" ] || [ -z "$RECOVERY_TYPE" ]; then
    print_error "Environment and recovery type are required"
    show_usage
    exit 1
fi

# Set environment-specific variables
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        STORAGE_CLASS="cnpg-local-storage"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        STORAGE_CLASS="cnpg-local-storage"
        ;;
    production)
        NAMESPACE="naytife-production"
        STORAGE_CLASS="cnpg-oci-storage"
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        exit 1
        ;;
esac

# Point-in-time recovery
restore_to_point_in_time() {
    if [ -z "$TARGET_TIME" ]; then
        print_error "Target time is required for PITR"
        exit 1
    fi
    
    print_header "Point-in-time recovery to $TARGET_TIME"
    
    RECOVERY_MANIFEST=$(cat <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: $NEW_CLUSTER
  namespace: $NAMESPACE
  labels:
    app.kubernetes.io/name: $NEW_CLUSTER
    app.kubernetes.io/managed-by: cnpg-recovery
    app.kubernetes.io/component: cnpg-cluster
    app.kubernetes.io/part-of: naytife-platform
spec:
  instances: 1
  
  bootstrap:
    recovery:
      source: naytife-postgres
      recoveryTarget:
        targetTime: "$TARGET_TIME"
      
  storage:
    size: 50Gi
    storageClass: $STORAGE_CLASS
    
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
      
  monitoring:
    enabled: true
EOF
)
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would create recovery cluster with manifest:"
        echo "$RECOVERY_MANIFEST"
    else
        echo "$RECOVERY_MANIFEST" | kubectl apply -f -
        print_success "Recovery cluster created: $NEW_CLUSTER"
    fi
}

# Full backup restore
restore_from_backup() {
    if [ -z "$BACKUP_NAME" ]; then
        print_error "Backup name is required for backup restore"
        exit 1
    fi
    
    print_header "Restoring from backup: $BACKUP_NAME"
    
    RECOVERY_MANIFEST=$(cat <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: $NEW_CLUSTER
  namespace: $NAMESPACE
  labels:
    app.kubernetes.io/name: $NEW_CLUSTER
    app.kubernetes.io/managed-by: cnpg-recovery
    app.kubernetes.io/component: cnpg-cluster
    app.kubernetes.io/part-of: naytife-platform
spec:
  instances: 1
  
  bootstrap:
    recovery:
      backup:
        name: "$BACKUP_NAME"
      
  storage:
    size: 50Gi
    storageClass: $STORAGE_CLASS
    
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
      
  monitoring:
    enabled: true
EOF
)
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would create recovery cluster with manifest:"
        echo "$RECOVERY_MANIFEST"
    else
        echo "$RECOVERY_MANIFEST" | kubectl apply -f -
        print_success "Recovery cluster created: $NEW_CLUSTER"
    fi
}

# Clone from existing cluster
clone_cluster() {
    if [ -z "$SOURCE_CLUSTER" ]; then
        SOURCE_CLUSTER="naytife-postgres"
    fi
    
    print_header "Cloning from cluster: $SOURCE_CLUSTER"
    
    RECOVERY_MANIFEST=$(cat <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: $NEW_CLUSTER
  namespace: $NAMESPACE
  labels:
    app.kubernetes.io/name: $NEW_CLUSTER
    app.kubernetes.io/managed-by: cnpg-recovery
    app.kubernetes.io/component: cnpg-cluster
    app.kubernetes.io/part-of: naytife-platform
spec:
  instances: 1
  
  bootstrap:
    recovery:
      source: $SOURCE_CLUSTER
      
  storage:
    size: 50Gi
    storageClass: $STORAGE_CLASS
    
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
      
  monitoring:
    enabled: true
EOF
)
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would create cloned cluster with manifest:"
        echo "$RECOVERY_MANIFEST"
    else
        echo "$RECOVERY_MANIFEST" | kubectl apply -f -
        print_success "Cloned cluster created: $NEW_CLUSTER"
    fi
}

# Validate recovery
validate_recovery() {
    if [ "$DRY_RUN" = true ]; then
        print_info "Would validate recovery cluster: $NEW_CLUSTER"
        return
    fi
    
    print_header "Validating recovery cluster"
    
    # Wait for cluster to be ready
    print_info "Waiting for recovery cluster to be ready..."
    kubectl wait --for=condition=ClusterReady cluster/"$NEW_CLUSTER" -n "$NAMESPACE" --timeout=600s
    
    # Test database connectivity
    RECOVERY_PRIMARY=$(kubectl get pods -n "$NAMESPACE" -l postgresql="$NEW_CLUSTER",role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$RECOVERY_PRIMARY" ]; then
        print_info "Testing database connectivity..."
        if kubectl exec -n "$NAMESPACE" "$RECOVERY_PRIMARY" -- psql -U naytife -d naytifedb -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema IN ('hydra', 'naytife_schema');" >/dev/null 2>&1; then
            print_success "Database connectivity test passed"
        else
            print_error "Database connectivity test failed"
            exit 1
        fi
    else
        print_error "Cannot find recovery primary pod for testing"
        exit 1
    fi
    
    print_success "Recovery validation passed"
}

# List available backups
list_backups() {
    print_header "Available backups"
    
    if kubectl get backups -n "$NAMESPACE" >/dev/null 2>&1; then
        kubectl get backups -n "$NAMESPACE" -o custom-columns=NAME:.metadata.name,STARTED:.status.startedAt,COMPLETED:.status.completedAt,STATUS:.status.phase
    else
        print_info "No backups found in namespace $NAMESPACE"
    fi
}

# Main execution
main() {
    print_header "CNPG Recovery for $ENVIRONMENT environment"
    echo "============================================="
    
    # Check if kubectl is working
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    # Check if CNPG operator is installed
    if ! kubectl get crd clusters.postgresql.cnpg.io >/dev/null 2>&1; then
        print_error "CNPG operator not installed"
        exit 1
    fi
    
    # Check if recovery cluster already exists
    if kubectl get cluster "$NEW_CLUSTER" -n "$NAMESPACE" >/dev/null 2>&1; then
        print_error "Recovery cluster '$NEW_CLUSTER' already exists"
        print_info "Delete it first: kubectl delete cluster '$NEW_CLUSTER' -n '$NAMESPACE'"
        exit 1
    fi
    
    case $RECOVERY_TYPE in
        pitr)
            restore_to_point_in_time
            ;;
        backup)
            restore_from_backup
            ;;
        clone)
            clone_cluster
            ;;
        *)
            print_error "Invalid recovery type: $RECOVERY_TYPE"
            exit 1
            ;;
    esac
    
    validate_recovery
    
    echo ""
    print_success "Recovery completed successfully!"
    echo ""
    print_info "Recovery cluster: $NEW_CLUSTER"
    print_info "Namespace: $NAMESPACE"
    print_info ""
    print_info "Next steps:"
    print_info "  1. Verify data integrity in recovered cluster"
    print_info "  2. Test application connectivity"
    print_info "  3. If satisfied, promote recovered cluster to primary"
    print_info "  4. Update application configuration"
    print_info "  5. Remove old cluster when stable"
    print_info ""
    print_info "Useful commands:"
    print_info "  • Check cluster: kubectl get cluster $NEW_CLUSTER -n $NAMESPACE"
    print_info "  • Connect to DB: kubectl exec -it \$(kubectl get pods -n $NAMESPACE -l postgresql=$NEW_CLUSTER,role=primary -o jsonpath='{.items[0].metadata.name}') -- psql -U naytife naytifedb"
    print_info "  • Delete recovery cluster: kubectl delete cluster $NEW_CLUSTER -n $NAMESPACE"
}

# Handle special case for listing backups
if [ "$1" = "list-backups" ]; then
    if [ -z "$2" ]; then
        print_error "Environment is required for listing backups"
        show_usage
        exit 1
    fi
    ENVIRONMENT="$2"
    case $ENVIRONMENT in
        local) NAMESPACE="naytife" ;;
        staging) NAMESPACE="naytife-staging" ;;
        production) NAMESPACE="naytife-production" ;;
        *) print_error "Invalid environment: $ENVIRONMENT"; exit 1 ;;
    esac
    list_backups
    exit 0
fi

# Run main function
main "$@"
