#!/bin/bash

# Cleanup Script for Kustomize-based deployments
# Usage: ./cleanup.sh [environment] [--force] [--dry-run]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}🧹 $1${NC}"; }

# Default values
ENVIRONMENT=""
FORCE=false
DRY_RUN=false

show_usage() {
    echo -e "${CYAN}🧹 Naytife Cleanup Tool${NC}"
    echo "======================="
    echo "Usage: $0 <environment> [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --force       Skip confirmation prompts"
    echo "  --dry-run     Show what would be deleted without actually deleting"
    echo "  --help, -h    Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                      # Clean up local environment"
    echo "  $0 staging --dry-run          # Preview staging cleanup"
    echo "  $0 production --force         # Force cleanup production (dangerous!)"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --force)
            FORCE=true
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

# Check if environment is provided
if [ -z "$ENVIRONMENT" ]; then
    print_error "Environment is required"
    show_usage
    exit 1
fi

# Set namespaces based on environment
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        AUTH_NAMESPACE="naytife-auth"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        AUTH_NAMESPACE="naytife-auth-staging"
        ;;
    production)
        NAMESPACE="naytife-production"
        AUTH_NAMESPACE="naytife-auth-production"
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        exit 1
        ;;
esac

print_header "Cleanup for $ENVIRONMENT Environment"
echo "======================================="

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed"
    exit 1
fi

# Check cluster connectivity
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

# Safety check for production
if [ "$ENVIRONMENT" = "production" ] && [ "$FORCE" = false ]; then
    print_warning "You are about to clean up the PRODUCTION environment!"
    print_warning "This will delete ALL resources in the production namespace."
    echo ""
    read -p "Are you absolutely sure you want to continue? (type 'DELETE PRODUCTION' to confirm): " confirmation
    
    if [ "$confirmation" != "DELETE PRODUCTION" ]; then
        print_info "Cleanup cancelled"
        exit 0
    fi
fi

# Check if namespaces exist
MAIN_NS_EXISTS=false
AUTH_NS_EXISTS=false

if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    MAIN_NS_EXISTS=true
    print_info "Target namespace: $NAMESPACE"
else
    print_warning "Main namespace '$NAMESPACE' does not exist"
fi

if kubectl get namespace "$AUTH_NAMESPACE" >/dev/null 2>&1; then
    AUTH_NS_EXISTS=true
    print_info "Auth namespace: $AUTH_NAMESPACE"
else
    print_warning "Auth namespace '$AUTH_NAMESPACE' does not exist"
fi

if [ "$MAIN_NS_EXISTS" = false ] && [ "$AUTH_NS_EXISTS" = false ]; then
    print_warning "No target namespaces exist"
    exit 0
fi

# Function to show resources that would be deleted
show_resources() {
    local namespace=$1
    
    print_header "Resources in namespace $namespace:"
    
    echo ""
    print_info "Deployments:"
    kubectl get deployments -n "$namespace" 2>/dev/null || print_info "No deployments found"
    
    echo ""
    print_info "Services:"
    kubectl get services -n "$namespace" 2>/dev/null || print_info "No services found"
    
    echo ""
    print_info "ConfigMaps:"
    kubectl get configmaps -n "$namespace" 2>/dev/null || print_info "No configmaps found"
    
    echo ""
    print_info "Secrets:"
    kubectl get secrets -n "$namespace" 2>/dev/null || print_info "No secrets found"
    
    echo ""
    print_info "Persistent Volume Claims:"
    kubectl get pvc -n "$namespace" 2>/dev/null || print_info "No PVCs found"
    
    echo ""
    print_info "Ingress:"
    kubectl get ingress -n "$namespace" 2>/dev/null || print_info "No ingress resources found"
    
    echo ""
    print_info "Pods:"
    kubectl get pods -n "$namespace" 2>/dev/null || print_info "No pods found"
}

# Show what will be deleted
if [ "$MAIN_NS_EXISTS" = true ]; then
    show_resources "$NAMESPACE"
fi

if [ "$AUTH_NS_EXISTS" = true ]; then
    echo ""
    show_resources "$AUTH_NAMESPACE"
fi

if [ "$DRY_RUN" = true ]; then
    print_success "Dry run completed - no resources were deleted"
    exit 0
fi

# Confirmation for non-production environments
if [ "$FORCE" = false ]; then
    echo ""
    print_warning "This will delete ALL resources shown above"
    if [ "$MAIN_NS_EXISTS" = true ]; then
        print_warning "  - in namespace '$NAMESPACE'"
    fi
    if [ "$AUTH_NS_EXISTS" = true ]; then
        print_warning "  - in namespace '$AUTH_NAMESPACE'"
    fi
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Cleanup cancelled"
        exit 0
    fi
fi

# Function to delete resources with error handling
delete_resources() {
    local resource_type=$1
    local namespace=$2
    
    print_info "Deleting $resource_type..."
    
    if kubectl get "$resource_type" -n "$namespace" >/dev/null 2>&1; then
        local count=$(kubectl get "$resource_type" -n "$namespace" --no-headers 2>/dev/null | wc -l)
        if [ "$count" -gt 0 ]; then
            kubectl delete "$resource_type" --all -n "$namespace" --timeout=60s || print_warning "Some $resource_type may not have been deleted"
            print_success "Deleted $count $resource_type"
        else
            print_info "No $resource_type to delete"
        fi
    else
        print_info "No $resource_type found"
    fi
}

# Function to cleanup a single namespace
cleanup_namespace() {
    local ns_name=$1
    local ns_description="$2"
    
    if ! kubectl get namespace "$ns_name" >/dev/null 2>&1; then
        print_info "Namespace '$ns_name' doesn't exist, skipping"
        return
    fi
    
    print_header "Cleaning up $ns_description ($ns_name)"
    
    # Delete resources in order (least dependent first)
    print_info "Deleting application resources..."

    # Delete ingress first to stop external traffic
    delete_resources "ingress" "$ns_name"

    # Delete services to stop internal traffic
    delete_resources "services" "$ns_name"

    # Delete deployments
    delete_resources "deployments" "$ns_name"

    # Delete statefulsets if any
    delete_resources "statefulsets" "$ns_name"

    # Delete daemonsets if any
    delete_resources "daemonsets" "$ns_name"

    # Delete jobs and cronjobs
    delete_resources "jobs" "$ns_name"
    delete_resources "cronjobs" "$ns_name"

    # Delete pods (should be cleaned up by deployments, but just in case)
    print_info "Checking for remaining pods..."
    remaining_pods=$(kubectl get pods -n "$ns_name" --no-headers 2>/dev/null | wc -l)
    if [ "$remaining_pods" -gt 0 ]; then
        print_warning "Found $remaining_pods remaining pods, force deleting..."
        kubectl delete pods --all -n "$ns_name" --force --grace-period=0 || true
    fi

    # Delete configuration resources
    print_info "Deleting configuration resources..."
    delete_resources "configmaps" "$ns_name"
    delete_resources "secrets" "$ns_name"

    # Delete storage resources
    print_info "Deleting storage resources..."
    delete_resources "pvc" "$ns_name"

    # Delete RBAC resources
    delete_resources "serviceaccounts" "$ns_name"
    delete_resources "roles" "$ns_name"
    delete_resources "rolebindings" "$ns_name"

    # Wait for pods to terminate
    print_info "Waiting for all pods to terminate..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        remaining=$(kubectl get pods -n "$ns_name" --no-headers 2>/dev/null | wc -l)
        if [ "$remaining" -eq 0 ]; then
            break
        fi
        echo "Waiting for $remaining pods to terminate... (${timeout}s remaining)"
        sleep 5
        timeout=$((timeout - 5))
    done

    # Final check
    remaining_pods=$(kubectl get pods -n "$ns_name" --no-headers 2>/dev/null | wc -l)
    if [ "$remaining_pods" -gt 0 ]; then
        print_warning "$remaining_pods pods are still terminating in $ns_name"
        kubectl get pods -n "$ns_name"
    else
        print_success "All pods have been terminated in $ns_name"
    fi
    
    echo ""
}

print_header "Starting cleanup process..."

# Cleanup auth namespace first (dependencies might flow from main to auth)
if [ "$AUTH_NS_EXISTS" = true ]; then
    cleanup_namespace "$AUTH_NAMESPACE" "auth services"
fi

# Then cleanup main namespace
if [ "$MAIN_NS_EXISTS" = true ]; then
    cleanup_namespace "$NAMESPACE" "core services"
fi

# Option to delete the namespaces themselves
if [ "$ENVIRONMENT" = "local" ]; then
    if [ "$FORCE" = false ]; then
        if [ "$MAIN_NS_EXISTS" = true ]; then
            read -p "Do you want to delete the namespace '$NAMESPACE' itself? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                print_info "Deleting namespace $NAMESPACE..."
                kubectl delete namespace "$NAMESPACE" || print_warning "Failed to delete namespace"
            fi
        fi
        
        if [ "$AUTH_NS_EXISTS" = true ]; then
            read -p "Do you want to delete the auth namespace '$AUTH_NAMESPACE' itself? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                print_info "Deleting namespace $AUTH_NAMESPACE..."
                kubectl delete namespace "$AUTH_NAMESPACE" || print_warning "Failed to delete namespace"
            fi
        fi
    else
        if [ "$MAIN_NS_EXISTS" = true ]; then
            print_info "Deleting namespace $NAMESPACE..."
            kubectl delete namespace "$NAMESPACE" || print_warning "Failed to delete namespace"
        fi
        if [ "$AUTH_NS_EXISTS" = true ]; then
            print_info "Deleting namespace $AUTH_NAMESPACE..."
            kubectl delete namespace "$AUTH_NAMESPACE" || print_warning "Failed to delete namespace"
        fi
    fi
else
    print_info "Keeping namespaces for $ENVIRONMENT environment"
fi

echo ""
print_success "Cleanup completed for $ENVIRONMENT environment"

# Environment-specific cleanup notes
case $ENVIRONMENT in
    local)
        print_info "Local cleanup notes:"
        print_info "- Local k3s cluster is still running"
        print_info "- To completely reset, consider recreating the cluster"
        ;;
    staging|production)
        print_info "$ENVIRONMENT cleanup notes:"
        print_info "- Namespace and cluster are preserved"
        print_info "- DNS records and external resources unchanged"
        print_info "- Re-deploy when ready using: ./deploy.sh $ENVIRONMENT"
        ;;
esac
