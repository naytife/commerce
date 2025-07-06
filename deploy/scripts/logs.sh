#!/bin/bash

# Enhanced Log Aggregation Script for Kustomize-based deployments
# Usage: ./logs.sh [service] [environment] [--follow] [--tail=N]

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
print_header() { echo -e "${CYAN}📋 $1${NC}"; }

# Default values
SERVICE=""
ENVIRONMENT="local"
FOLLOW=false
TAIL_LINES=100

# Service mappings function
get_service_name() {
    local service=$1
    case $service in
        backend) echo "backend" ;;
        auth-handler) echo "auth-handler" ;;
        postgres) echo "postgres" ;;
        redis) echo "redis" ;;
        hydra) echo "hydra" ;;
        oathkeeper) echo "oathkeeper" ;;
        store-deployer) echo "store-deployer" ;;
        template-registry) echo "template-registry" ;;
        *) echo "$service" ;;
    esac
}

# Function to get the correct namespace for a service
get_service_namespace() {
    local service=$1
    local environment=$2
    
    case $service in
        auth-handler|hydra|oathkeeper)
            case $environment in
                local) echo "naytife-auth" ;;
                staging) echo "naytife-auth-staging" ;;
                production) echo "naytife-auth-production" ;;
            esac
            ;;
        *)
            case $environment in
                local) echo "naytife" ;;
                staging) echo "naytife-staging" ;;
                production) echo "naytife-production" ;;
            esac
            ;;
    esac
}

show_usage() {
    echo -e "${CYAN}📋 Naytife Log Viewer${NC}"
    echo "===================="
    echo "Usage: $0 [service] [environment] [options]"
    echo ""
    echo -e "${YELLOW}Available services:${NC}"
    echo "  • backend"
    echo "  • auth-handler" 
    echo "  • postgres"
    echo "  • redis"
    echo "  • hydra"
    echo "  • oathkeeper"
    echo "  • store-deployer"
    echo "  • template-registry"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --follow, -f       Follow log output"
    echo "  --tail=N           Show last N lines (default: 100)"
    echo "  --help, -h         Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 backend                    # View backend logs (local env)"
    echo "  $0 backend staging            # View backend logs (staging env)"
    echo "  $0 backend local --follow     # Follow backend logs"
    echo "  $0 postgres local --tail=50   # Show last 50 lines of postgres logs"
    echo "  $0                            # Show all service logs"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        backend|auth-handler|postgres|redis|hydra|oathkeeper|store-deployer|template-registry)
            SERVICE="$1"
            shift
            ;;
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --follow|-f)
            FOLLOW=true
            shift
            ;;
        --tail=*)
            TAIL_LINES="${1#*=}"
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

# Set namespace based on environment
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        ;;
    production)
        NAMESPACE="naytife-production"
        ;;
esac

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed"
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

# Function to get pods for a service
get_service_pods() {
    local service=$1
    local namespace=$2
    kubectl get pods -n "$namespace" -l "app=$service" --no-headers -o custom-columns="NAME:.metadata.name" 2>/dev/null
}

# Function to show logs for a specific service
show_service_logs() {
    local service=$1
    local namespace=$2
    local follow_flag=""
    local tail_flag="--tail=$TAIL_LINES"
    
    if [ "$FOLLOW" = true ]; then
        follow_flag="-f"
    fi
    
    print_header "Logs for $service in $ENVIRONMENT environment"
    
    local pods=$(get_service_pods "$service" "$namespace")
    
    if [ -z "$pods" ]; then
        print_warning "No pods found for service '$service' in namespace '$namespace'"
        return 1
    fi
    
    echo "Found pods:"
    echo "$pods" | while read pod; do
        echo "  • $pod"
    done
    echo ""
    
    # If multiple pods, show logs from all
    local pod_count=$(echo "$pods" | wc -l)
    
    if [ "$pod_count" -eq 1 ]; then
        local pod=$(echo "$pods" | head -n1)
        print_info "Showing logs from pod: $pod"
        kubectl logs $follow_flag $tail_flag "$pod" -n "$namespace"
    else
        print_info "Multiple pods found. Showing logs from all pods:"
        echo "$pods" | while read pod; do
            echo ""
            print_info "=== Logs from $pod ==="
            kubectl logs $tail_flag "$pod" -n "$namespace" | sed "s/^/[$pod] /"
        done
        
        if [ "$FOLLOW" = true ]; then
            print_info "Following logs from all pods (multiplexed):"
            echo "$pods" | while read pod; do
                kubectl logs -f "$pod" -n "$namespace" | sed "s/^/[$pod] /" &
            done
            wait
        fi
    fi
}

# Function to show all logs
show_all_logs() {
    local namespace=$1
    
    print_header "All Service Logs in $ENVIRONMENT environment"
    
    # Get all deployments
    local deployments=$(kubectl get deployments -n "$namespace" --no-headers -o custom-columns="NAME:.metadata.name" 2>/dev/null)
    
    if [ -z "$deployments" ]; then
        print_warning "No deployments found in namespace '$namespace'"
        return 1
    fi
    
    echo "Available services:"
    echo "$deployments" | while read deployment; do
        echo "  • $deployment"
    done
    echo ""
    
    # Show recent logs from each service
    echo "$deployments" | while read deployment; do
        echo ""
        print_info "=== Recent logs from $deployment ==="
        local pods=$(get_service_pods "$deployment" "$namespace")
        if [ -n "$pods" ]; then
            echo "$pods" | head -n1 | while read pod; do
                kubectl logs --tail=20 "$pod" -n "$namespace" 2>/dev/null | sed "s/^/[$deployment] /" || print_warning "Could not get logs for $pod"
            done
        else
            print_warning "No pods found for $deployment"
        fi
    done
}

# Main logic
if [ -z "$SERVICE" ]; then
    show_all_logs "$NAMESPACE"
else
    if [[ ! "$SERVICE" =~ ^(backend|auth-handler|postgres|redis|hydra|oathkeeper|store-deployer|template-registry)$ ]]; then
        print_error "Unknown service: $SERVICE"
        show_usage
        exit 1
    fi
    
    # Get the correct namespace for this service
    SERVICE_NAMESPACE=$(get_service_namespace "$SERVICE" "$ENVIRONMENT")
    show_service_logs "$SERVICE" "$SERVICE_NAMESPACE"
fi

print_success "Log viewing completed"
