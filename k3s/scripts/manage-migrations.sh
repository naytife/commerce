#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

NAMESPACE="naytife"
JOB_NAME="backend-migrate"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

show_help() {
    echo -e "${BLUE}Migration Management Script${NC}"
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  status           - Show current migration status"
    echo "  logs             - Show migration job logs"
    echo "  run              - Run migrations manually"
    echo "  dry-run          - Test migrations without applying"
    echo "  rollback [VER]   - Rollback to specific version"
    echo "  clean            - Clean failed migration jobs"
    echo "  validate         - Validate migration files"
    echo "  help             - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status"
    echo "  $0 run"
    echo "  $0 rollback 20250603120000"
    echo "  $0 logs"
}

check_prerequisites() {
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is required but not installed${NC}"
        exit 1
    fi
    
    if ! kubectl get namespace $NAMESPACE &> /dev/null; then
        echo -e "${RED}❌ Namespace $NAMESPACE not found${NC}"
        exit 1
    fi
}

get_migration_status() {
    echo -e "${BLUE}📊 Migration Status${NC}"
    echo "==================="
    
    # Check if migration job exists
    if kubectl get job $JOB_NAME -n $NAMESPACE &> /dev/null; then
        JOB_STATUS=$(kubectl get job $JOB_NAME -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' 2>/dev/null || echo "Unknown")
        JOB_FAILED=$(kubectl get job $JOB_NAME -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Failed")].status}' 2>/dev/null || echo "False")
        
        if [[ "$JOB_STATUS" == "True" ]]; then
            echo -e "${GREEN}✅ Last migration job: Completed successfully${NC}"
        elif [[ "$JOB_FAILED" == "True" ]]; then
            echo -e "${RED}❌ Last migration job: Failed${NC}"
        else
            echo -e "${YELLOW}⏳ Last migration job: Running or pending${NC}"
        fi
        
        # Show job details
        echo ""
        kubectl get job $JOB_NAME -n $NAMESPACE
        echo ""
        kubectl get pods -l job-name=$JOB_NAME -n $NAMESPACE
    else
        echo -e "${YELLOW}⚠️  No migration job found${NC}"
    fi
}

show_logs() {
    echo -e "${BLUE}📋 Migration Logs${NC}"
    echo "=================="
    
    if kubectl get job $JOB_NAME -n $NAMESPACE &> /dev/null; then
        kubectl logs job/$JOB_NAME -n $NAMESPACE --tail=100 -f
    else
        echo -e "${RED}❌ No migration job found${NC}"
        exit 1
    fi
}

run_migration() {
    echo -e "${YELLOW}🚀 Running database migrations...${NC}"
    
    # Clean up any existing job
    kubectl delete job $JOB_NAME -n $NAMESPACE 2>/dev/null || true
    
    # Update ConfigMap with latest migration files
    BACKEND_DIR="$(cd "$SCRIPT_DIR/../../backend" && pwd)"
    MIGRATIONS_DIR="$BACKEND_DIR/internal/db/migrations"
    
    if [ -d "$MIGRATIONS_DIR" ]; then
        echo -e "${YELLOW}📝 Updating migration ConfigMap...${NC}"
        kubectl delete configmap backend-migrations -n $NAMESPACE 2>/dev/null || true
        kubectl create configmap backend-migrations -n $NAMESPACE \
            --from-file="$MIGRATIONS_DIR"
        echo -e "${GREEN}✅ Migration ConfigMap updated${NC}"
    else
        echo -e "${RED}❌ Migration directory not found at $MIGRATIONS_DIR${NC}"
        exit 1
    fi
    
    # Apply migration job
    MANIFESTS_DIR="$(cd "$SCRIPT_DIR/../manifests" && pwd)"
    kubectl apply -f "$MANIFESTS_DIR/06-backend/backend-migration.yaml"
    
    # Wait for completion
    echo -e "${YELLOW}⏳ Waiting for migration to complete...${NC}"
    kubectl wait --for=condition=complete job/$JOB_NAME -n $NAMESPACE --timeout=300s
    
    # Check result
    if kubectl get job $JOB_NAME -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' | grep -q "True"; then
        echo -e "${GREEN}✅ Migration completed successfully${NC}"
    else
        echo -e "${RED}❌ Migration failed${NC}"
        kubectl logs job/$JOB_NAME -n $NAMESPACE --tail=50
        exit 1
    fi
}

dry_run_migration() {
    echo -e "${YELLOW}🧪 Running migration dry-run...${NC}"
    
    # Create a temporary job for dry-run
    kubectl run migration-dry-run --rm -i --restart=Never --image=arigaio/atlas:latest -n $NAMESPACE -- \
        sh -c "atlas migrate apply --dir file:///migrations --url \$DATABASE_URL --dry-run" \
        --env DATABASE_URL="$(kubectl get secret backend-secret -n $NAMESPACE -o jsonpath='{.data.DATABASE_URL}' | base64 -d)"
}

rollback_migration() {
    local target_version=$1
    
    if [ -z "$target_version" ]; then
        echo -e "${RED}❌ Target version is required for rollback${NC}"
        echo "Usage: $0 rollback <version>"
        exit 1
    fi
    
    echo -e "${YELLOW}⚠️  Rolling back to version: $target_version${NC}"
    
    # Warning prompt
    read -p "Are you sure you want to rollback? This action cannot be undone. (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Rollback cancelled"
        exit 0
    fi
    
    # Create rollback job
    kubectl run migration-rollback --rm -i --restart=Never --image=arigaio/atlas:latest -n $NAMESPACE -- \
        sh -c "atlas migrate down --dir file:///migrations --url \$DATABASE_URL --to-version $target_version" \
        --env DATABASE_URL="$(kubectl get secret backend-secret -n $NAMESPACE -o jsonpath='{.data.DATABASE_URL}' | base64 -d)"
}

clean_jobs() {
    echo -e "${YELLOW}🧹 Cleaning migration jobs...${NC}"
    
    kubectl delete job $JOB_NAME -n $NAMESPACE 2>/dev/null || true
    kubectl delete pod -l job-name=$JOB_NAME -n $NAMESPACE 2>/dev/null || true
    
    echo -e "${GREEN}✅ Migration jobs cleaned${NC}"
}

validate_migrations() {
    echo -e "${BLUE}🔍 Validating migration files...${NC}"
    
    BACKEND_DIR="$(cd "$SCRIPT_DIR/../../backend" && pwd)"
    MIGRATIONS_DIR="$BACKEND_DIR/internal/db/migrations"
    
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        echo -e "${RED}❌ Migration directory not found${NC}"
        exit 1
    fi
    
    echo "Migration files found:"
    ls -la "$MIGRATIONS_DIR"
    
    # Use atlas to validate
    cd "$BACKEND_DIR"
    if command -v atlas &> /dev/null; then
        echo -e "${YELLOW}🔍 Running Atlas validation...${NC}"
        atlas migrate validate --env local
        echo -e "${GREEN}✅ Migration files are valid${NC}"
    else
        echo -e "${YELLOW}⚠️  Atlas not found, skipping validation${NC}"
    fi
}

# Main script logic
case "${1:-help}" in
    "status")
        check_prerequisites
        get_migration_status
        ;;
    "logs")
        check_prerequisites
        show_logs
        ;;
    "run")
        check_prerequisites
        run_migration
        ;;
    "dry-run")
        check_prerequisites
        dry_run_migration
        ;;
    "rollback")
        check_prerequisites
        rollback_migration "$2"
        ;;
    "clean")
        check_prerequisites
        clean_jobs
        ;;
    "validate")
        validate_migrations
        ;;
    "help"|*)
        show_help
        ;;
esac
