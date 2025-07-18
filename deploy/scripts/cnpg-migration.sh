#!/bin/bash

# CNPG Migration Script for Naytife Commerce Platform
# This script migrates from the current PostgreSQL deployment to CloudNativePG

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
    echo -e "${CYAN}📋 $1${NC}"
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
    echo -e "${CYAN}📦 CNPG Migration Script${NC}"
    echo "========================="
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --dry-run              Show what would be done without executing"
    echo "  --skip-backup          Skip creating backup (dangerous!)"
    echo "  --force                Force migration even if CNPG cluster exists"
    echo "  --help, -h             Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                     # Migrate local environment"
    echo "  $0 staging --dry-run         # Preview staging migration"
    echo "  $0 production --force        # Force production migration"
}

# Default values
ENVIRONMENT=""
DRY_RUN=false
SKIP_BACKUP=false
FORCE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --skip-backup)
            SKIP_BACKUP=true
            shift
            ;;
        --force)
            FORCE=true
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

# Check if environment is provided
if [ -z "$ENVIRONMENT" ]; then
    print_error "Environment is required"
    show_usage
    exit 1
fi

# Set environment-specific variables
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        CONTEXT_PATTERN="k3d"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        CONTEXT_PATTERN="staging"
        ;;
    production)
        NAMESPACE="naytife-production"
        CONTEXT_PATTERN="production"
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        exit 1
        ;;
esac

# Production safety check
if [ "$ENVIRONMENT" = "production" ] && [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
    print_warning "🚨 PRODUCTION MIGRATION WARNING 🚨"
    print_warning "You are about to migrate the production PostgreSQL database to CNPG!"
    echo ""
    print_info "Migration details:"
    print_info "  • Environment: $ENVIRONMENT"
    print_info "  • Namespace: $NAMESPACE"
    print_info "  • Source: current postgres deployment"
    print_info "  • Target: CloudNativePG cluster"
    echo ""
    read -p "Are you sure you want to continue? (type 'MIGRATE TO CNPG' to confirm): " confirmation
    
    if [ "$confirmation" != "MIGRATE TO CNPG" ]; then
        print_info "Production migration cancelled"
        exit 0
    fi
fi

print_header "CNPG Migration for $ENVIRONMENT environment"

# Pre-migration validation
validate_environment() {
    print_header "Pre-migration validation"
    
    # Check if kubectl is working
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    # Check if CNPG operator is installed
    if ! kubectl get crd clusters.postgresql.cnpg.io >/dev/null 2>&1; then
        print_error "CNPG operator not installed. Please install it first:"
        print_info "  kubectl apply -k deploy/base/cnpg-operator"
        exit 1
    fi
    
    # Check if current postgres deployment exists
    if ! kubectl get deployment postgres -n "$NAMESPACE" >/dev/null 2>&1; then
        print_error "Current postgres deployment not found in namespace $NAMESPACE"
        exit 1
    fi
    
    # Check if CNPG cluster already exists
    if kubectl get cluster naytife-postgres -n "$NAMESPACE" >/dev/null 2>&1; then
        if [ "$FORCE" = false ]; then
            print_error "CNPG cluster already exists. Use --force to override"
            exit 1
        else
            print_warning "CNPG cluster exists but --force specified, continuing..."
        fi
    fi
    
    print_success "Environment validation passed"
}

# Create database backup
create_backup() {
    if [ "$SKIP_BACKUP" = true ]; then
        print_warning "Skipping backup creation (--skip-backup specified)"
        return
    fi
    
    print_header "Creating database backup"
    
    # Get current postgres pod
    POSTGRES_POD=$(kubectl get pods -n "$NAMESPACE" -l app=postgres -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -z "$POSTGRES_POD" ]; then
        print_error "No postgres pod found"
        exit 1
    fi
    
    # Check if pod is ready
    if ! kubectl get pod "$POSTGRES_POD" -n "$NAMESPACE" -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' | grep -q "True"; then
        print_error "Postgres pod is not ready"
        exit 1
    fi
    
    BACKUP_FILE="/tmp/naytife-backup-$(date +%Y%m%d_%H%M%S).sql"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would create backup: $BACKUP_FILE"
    else
        print_info "Creating backup: $BACKUP_FILE"
        kubectl exec -n "$NAMESPACE" "$POSTGRES_POD" -- pg_dump -U naytife naytifedb > "$BACKUP_FILE"
        print_success "Backup created: $BACKUP_FILE"
    fi
}

# Deploy CNPG cluster
deploy_cnpg() {
    print_header "Deploying CNPG cluster"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would deploy CNPG cluster using:"
        print_info "  kubectl apply -k deploy/overlays/$ENVIRONMENT"
        return
    fi
    
    # Deploy CNPG operator first if not already deployed
    if ! kubectl get deployment cnpg-controller-manager -n cnpg-system >/dev/null 2>&1; then
        print_info "Installing CNPG operator..."
        kubectl apply -k "$DEPLOY_DIR/base/cnpg-operator"
        kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=300s
    fi
    
    # Deploy storage classes
    kubectl apply -k "$DEPLOY_DIR/base/cnpg-storage"
    
    # Deploy CNPG cluster and related resources
    kubectl apply -k "$DEPLOY_DIR/overlays/$ENVIRONMENT"
    
    # Wait for cluster to be ready
    print_info "Waiting for CNPG cluster to be ready..."
    kubectl wait --for=condition=ClusterReady cluster/naytife-postgres -n "$NAMESPACE" --timeout=600s
    
    print_success "CNPG cluster deployed successfully"
}

# Migrate data
migrate_data() {
    print_header "Migrating data to CNPG"
    
    if [ "$SKIP_BACKUP" = true ]; then
        print_warning "Skipping data migration (no backup created)"
        return
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would migrate data from backup to CNPG cluster"
        return
    fi
    
    # Get CNPG primary pod
    print_info "Finding CNPG primary pod..."
    CNPG_PRIMARY=""
    for i in {1..30}; do
        CNPG_PRIMARY=$(kubectl get pods -n "$NAMESPACE" -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
        if [ -n "$CNPG_PRIMARY" ]; then
            break
        fi
        print_info "Waiting for CNPG primary pod... ($i/30)"
        sleep 10
    done
    
    if [ -z "$CNPG_PRIMARY" ]; then
        print_error "CNPG primary pod not found after 5 minutes"
        exit 1
    fi
    
    print_info "CNPG primary pod: $CNPG_PRIMARY"
    
    # Wait for pod to be ready
    kubectl wait --for=condition=Ready pod/"$CNPG_PRIMARY" -n "$NAMESPACE" --timeout=300s
    
    # Restore backup
    BACKUP_FILE="/tmp/naytife-backup-$(date +%Y%m%d_%H%M%S).sql"
    if [ -f "$BACKUP_FILE" ]; then
        print_info "Restoring backup to CNPG cluster..."
        kubectl exec -n "$NAMESPACE" "$CNPG_PRIMARY" -- psql -U naytife -d naytifedb < "$BACKUP_FILE"
        print_success "Data migration completed"
    else
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
}

# Update application configuration
update_app_config() {
    print_header "Updating application configuration"
    
    # Create new DATABASE_URL using pooler
    NEW_DATABASE_URL="postgresql://naytife:\${POSTGRES_PASSWORD}@naytife-postgres-pooler.${NAMESPACE}.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
    NEW_DATABASE_URL_DIRECT="postgresql://naytife:\${POSTGRES_PASSWORD}@naytife-postgres-rw.${NAMESPACE}.svc.cluster.local:5432/naytifedb?sslmode=require&search_path=naytife_schema"
    
    print_info "New connection strings:"
    print_info "  Pooled: $NEW_DATABASE_URL"
    print_info "  Direct: $NEW_DATABASE_URL_DIRECT"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would update backend-secret.yaml with new connection strings"
    else
        print_info "Manual step required: Update deploy/secrets/${ENVIRONMENT}/backend-secret.yaml"
        print_info "Replace DATABASE_URL with the pooled connection string"
        print_info "Add DATABASE_URL_DIRECT with the direct connection string for migrations"
    fi
    
    print_success "Configuration update guidance provided"
}

# Validation and cleanup
validate_migration() {
    print_header "Validating migration"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "Would validate CNPG cluster and database connectivity"
        return
    fi
    
    # Test CNPG cluster status
    if kubectl get cluster naytife-postgres -n "$NAMESPACE" -o jsonpath='{.status.phase}' | grep -q "Cluster in healthy state"; then
        print_success "CNPG cluster is healthy"
    else
        print_error "CNPG cluster is not healthy"
        exit 1
    fi
    
    # Test database connectivity
    CNPG_PRIMARY=$(kubectl get pods -n "$NAMESPACE" -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$CNPG_PRIMARY" ]; then
        print_info "Testing database connectivity..."
        if kubectl exec -n "$NAMESPACE" "$CNPG_PRIMARY" -- psql -U naytife -d naytifedb -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema IN ('hydra', 'naytife_schema');" >/dev/null 2>&1; then
            print_success "Database connectivity test passed"
        else
            print_error "Database connectivity test failed"
            exit 1
        fi
    else
        print_error "Cannot find CNPG primary pod for testing"
        exit 1
    fi
    
    print_success "Migration validation passed"
}

# Main execution
main() {
    print_header "CNPG Migration for $ENVIRONMENT environment"
    echo "================================================="
    
    validate_environment
    create_backup
    deploy_cnpg
    migrate_data
    update_app_config
    validate_migration
    
    echo ""
    print_success "Migration completed successfully!"
    echo ""
    print_info "Next steps:"
    print_info "  1. Update application secrets with new connection strings"
    print_info "  2. Restart application deployments"
    print_info "  3. Verify application connectivity"
    print_info "  4. Remove old postgres deployment when stable"
    print_info "  5. Monitor CNPG cluster health"
    
    if [ "$ENVIRONMENT" = "local" ]; then
        print_info ""
        print_info "Local development commands:"
        print_info "  • Check cluster: kubectl get cluster naytife-postgres -n $NAMESPACE"
        print_info "  • Check pooler: kubectl get pooler naytife-postgres-pooler -n $NAMESPACE"
        print_info "  • Connect to DB: kubectl exec -it naytife-postgres-1 -n $NAMESPACE -- psql -U naytife naytifedb"
    fi
}

# Run main function
main "$@"
