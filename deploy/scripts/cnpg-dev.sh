#!/bin/bash

# CNPG Development Helper Script
# This script provides quick commands for CNPG development and testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}🛠️ $1${NC}"
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
    echo -e "${CYAN}🛠️ CNPG Development Helper${NC}"
    echo "============================="
    echo "Usage: $0 [command] [environment] [options]"
    echo ""
    echo -e "${YELLOW}Commands:${NC}"
    echo "  status       Show CNPG cluster status"
    echo "  connect      Connect to database"
    echo "  logs         Show cluster logs"
    echo "  backup       Create manual backup"
    echo "  restore      List restore options"
    echo "  test         Run integration tests"
    echo "  metrics      Show metrics endpoint"
    echo "  cleanup      Clean up CNPG resources"
    echo "  install      Deploy CNPG operator and cluster (GitOps)"
    echo "  deploy       Deploy CNPG cluster"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 status local              # Show local cluster status"
    echo "  $0 connect staging           # Connect to staging database"
    echo "  $0 logs production           # Show production logs"
    echo "  $0 backup local              # Create local backup"
    echo "  $0 test local --verbose      # Run detailed tests"
}

# Default values
COMMAND=""
ENVIRONMENT="local"
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        status|connect|logs|backup|restore|test|metrics|cleanup|install|deploy)
            COMMAND="$1"
            shift
            ;;
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            if [ -z "$COMMAND" ]; then
                COMMAND="$1"
            fi
            shift
            ;;
    esac
done

# Check if command is provided
if [ -z "$COMMAND" ]; then
    show_usage
    exit 1
fi

# Set environment-specific variables
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        CLUSTER_NAME="naytife-postgres"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        CLUSTER_NAME="naytife-postgres"
        ;;
    production)
        NAMESPACE="naytife-production"
        CLUSTER_NAME="naytife-postgres"
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        exit 1
        ;;
esac

# Check prerequisites
check_prerequisites() {
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed"
        exit 1
    fi
    
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
}

# Show cluster status
show_status() {
    print_header "CNPG Cluster Status ($ENVIRONMENT)"
    
    # Check if cluster exists
    if ! kubectl get cluster $CLUSTER_NAME -n $NAMESPACE >/dev/null 2>&1; then
        print_error "CNPG cluster '$CLUSTER_NAME' not found in namespace '$NAMESPACE'"
        return 1
    fi
    
    # Cluster overview
    echo -e "${BLUE}Cluster Overview:${NC}"
    kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o wide
    
    echo ""
    echo -e "${BLUE}Pod Status:${NC}"
    kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME
    
    echo ""
    echo -e "${BLUE}Services:${NC}"
    kubectl get services -n $NAMESPACE | grep $CLUSTER_NAME
    
    # Pooler status
    if kubectl get pooler $CLUSTER_NAME-pooler -n $NAMESPACE >/dev/null 2>&1; then
        echo ""
        echo -e "${BLUE}Pooler Status:${NC}"
        kubectl get pooler $CLUSTER_NAME-pooler -n $NAMESPACE
    fi
    
    # PVC status
    echo ""
    echo -e "${BLUE}Storage:${NC}"
    kubectl get pvc -n $NAMESPACE | grep $CLUSTER_NAME
    
    # Recent events
    echo ""
    echo -e "${BLUE}Recent Events:${NC}"
    kubectl get events -n $NAMESPACE --sort-by=.metadata.creationTimestamp | grep -i postgres | tail -5
}

# Connect to database
connect_db() {
    print_header "Connecting to Database ($ENVIRONMENT)"
    
    # Get primary pod
    PRIMARY_POD=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -z "$PRIMARY_POD" ]; then
        print_error "No primary pod found"
        return 1
    fi
    
    print_info "Connecting to primary pod: $PRIMARY_POD"
    print_info "Use \\q to exit"
    
    kubectl exec -it $PRIMARY_POD -n $NAMESPACE -- psql -U naytife naytifedb
}

# Show logs
show_logs() {
    print_header "CNPG Logs ($ENVIRONMENT)"
    
    # Primary pod logs
    PRIMARY_POD=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -n "$PRIMARY_POD" ]; then
        echo -e "${BLUE}Primary Pod Logs:${NC}"
        kubectl logs $PRIMARY_POD -n $NAMESPACE --tail=50
    fi
    
    # Pooler logs
    POOLER_POD=$(kubectl get pods -n $NAMESPACE -l cnpg.io/poolerName=$CLUSTER_NAME-pooler -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -n "$POOLER_POD" ]; then
        echo ""
        echo -e "${BLUE}Pooler Logs:${NC}"
        kubectl logs $POOLER_POD -n $NAMESPACE --tail=50
    fi
}

# Create backup
create_backup() {
    print_header "Creating Manual Backup ($ENVIRONMENT)"
    
    # Check if cluster exists
    if ! kubectl get cluster $CLUSTER_NAME -n $NAMESPACE >/dev/null 2>&1; then
        print_error "CNPG cluster '$CLUSTER_NAME' not found"
        return 1
    fi
    
    # Create backup annotation
    kubectl annotate cluster $CLUSTER_NAME -n $NAMESPACE cnpg.io/backup=manual-$(date +%Y%m%d-%H%M%S)
    
    print_success "Manual backup initiated"
    print_info "Monitor backup progress with: kubectl get backups -n $NAMESPACE"
}

# Show restore options
show_restore() {
    print_header "Restore Options ($ENVIRONMENT)"
    
    echo -e "${BLUE}Available Backups:${NC}"
    kubectl get backups -n $NAMESPACE 2>/dev/null || echo "No backups found"
    
    echo ""
    echo -e "${BLUE}Restore Commands:${NC}"
    echo "Point-in-time recovery:"
    echo "  ./cnpg-recovery.sh $ENVIRONMENT pitr --target-time=2024-07-15T10:30:00Z"
    echo ""
    echo "Restore from backup:"
    echo "  ./cnpg-recovery.sh $ENVIRONMENT backup --backup-name=<backup-name>"
    echo ""
    echo "Clone cluster:"
    echo "  ./cnpg-recovery.sh $ENVIRONMENT clone --source-cluster=$CLUSTER_NAME"
}

# Run tests
run_tests() {
    print_header "Running CNPG Integration Tests ($ENVIRONMENT)"
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    if [ "$VERBOSE" = true ]; then
        $SCRIPT_DIR/test-cnpg-integration.sh $ENVIRONMENT --verbose
    else
        $SCRIPT_DIR/test-cnpg-integration.sh $ENVIRONMENT
    fi
}

# Show metrics
show_metrics() {
    print_header "CNPG Metrics ($ENVIRONMENT)"
    
    # Get primary pod
    PRIMARY_POD=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    
    if [ -z "$PRIMARY_POD" ]; then
        print_error "No primary pod found"
        return 1
    fi
    
    print_info "Metrics available at: http://localhost:9187/metrics"
    print_info "Starting port-forward... (Press Ctrl+C to stop)"
    
    kubectl port-forward $PRIMARY_POD 9187:9187 -n $NAMESPACE
}

# Clean up CNPG resources
cleanup_cnpg() {
    print_header "Cleaning up CNPG Resources ($ENVIRONMENT)"
    
    print_warning "This will delete the CNPG cluster and all data!"
    read -p "Are you sure? (type 'yes' to confirm): " confirmation
    
    if [ "$confirmation" != "yes" ]; then
        print_info "Cleanup cancelled"
        return 0
    fi
    
    # Delete cluster
    if kubectl get cluster $CLUSTER_NAME -n $NAMESPACE >/dev/null 2>&1; then
        kubectl delete cluster $CLUSTER_NAME -n $NAMESPACE
        print_success "CNPG cluster deleted"
    fi
    
    # Delete pooler
    if kubectl get pooler $CLUSTER_NAME-pooler -n $NAMESPACE >/dev/null 2>&1; then
        kubectl delete pooler $CLUSTER_NAME-pooler -n $NAMESPACE
        print_success "CNPG pooler deleted"
    fi
    
    # Delete backups
    kubectl delete backups -n $NAMESPACE -l cnpg.io/cluster=$CLUSTER_NAME 2>/dev/null || true
    
    print_success "CNPG cleanup completed"
}

# Install CNPG operator and cluster (GitOps approach)
install_operator() {
    print_header "Installing CNPG Operator and Cluster ($ENVIRONMENT)"
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Deploy complete CNPG setup using overlay (GitOps approach)
    kubectl apply -k $DEPLOY_DIR/overlays/$ENVIRONMENT
    
    print_info "Waiting for operator to be ready..."
    kubectl wait --for=condition=Available deployment/cnpg-controller-manager -n cnpg-system --timeout=300s
    
    print_info "Waiting for cluster to be ready..."
    kubectl wait --for=condition=ClusterReady cluster/naytife-postgres -n naytife --timeout=600s
    
    print_success "CNPG operator and cluster deployed successfully"
}

# Deploy CNPG cluster
deploy_cluster() {
    print_header "Deploying CNPG Cluster ($ENVIRONMENT)"
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Deploy storage classes
    kubectl apply -k $DEPLOY_DIR/base/cnpg-storage
    
    # Deploy cluster
    kubectl apply -k $DEPLOY_DIR/overlays/$ENVIRONMENT
    
    print_info "Waiting for cluster to be ready..."
    kubectl wait --for=condition=ClusterReady cluster/$CLUSTER_NAME -n $NAMESPACE --timeout=600s
    
    print_success "CNPG cluster deployed successfully"
}

# Main execution
main() {
    check_prerequisites
    
    case $COMMAND in
        status)
            show_status
            ;;
        connect)
            connect_db
            ;;
        logs)
            show_logs
            ;;
        backup)
            create_backup
            ;;
        restore)
            show_restore
            ;;
        test)
            run_tests
            ;;
        metrics)
            show_metrics
            ;;
        cleanup)
            cleanup_cnpg
            ;;
        install)
            install_operator
            ;;
        deploy)
            deploy_cluster
            ;;
        *)
            print_error "Unknown command: $COMMAND"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
