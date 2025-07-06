#!/bin/bash

# Debugging and Troubleshooting Script
# Usage: ./debug.sh [environment] [service] [--detailed]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}🔍 $1${NC}"; }
print_debug() { echo -e "${MAGENTA}🐛 $1${NC}"; }

# Default values
ENVIRONMENT="local"
SERVICE=""
DETAILED=false

show_usage() {
    echo -e "${CYAN}🐛 Naytife Debug Tool${NC}"
    echo "==================="
    echo "Usage: $0 [environment] [service] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo -e "${YELLOW}Services:${NC} backend, auth-handler, postgres, redis, hydra, oathkeeper"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --detailed    Show detailed debug information"
    echo "  --help, -h    Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                      # Debug all services in local"
    echo "  $0 local backend              # Debug backend service in local"
    echo "  $0 staging backend --detailed # Detailed debug of backend in staging"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        backend|auth-handler|postgres|redis|hydra|oathkeeper|store-deployer|template-registry)
            SERVICE="$1"
            shift
            ;;
        --detailed)
            DETAILED=true
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

print_header "Debug Information for $ENVIRONMENT Environment"
echo "================================================="

# Check basic connectivity
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

print_success "Connected to cluster: $(kubectl config current-context)"

# Function to debug a specific service
debug_service() {
    local service=$1
    local namespace=$2
    
    print_header "Debugging Service: $service"
    
    # Get deployment info
    echo ""
    print_debug "Deployment Status:"
    if kubectl get deployment "$service" -n "$namespace" >/dev/null 2>&1; then
        kubectl get deployment "$service" -n "$namespace" -o wide
        
        # Check deployment conditions
        echo ""
        print_debug "Deployment Conditions:"
        kubectl describe deployment "$service" -n "$namespace" | grep -A 10 "Conditions:"
        
        # Get replica set info
        echo ""
        print_debug "ReplicaSet Status:"
        kubectl get rs -n "$namespace" -l app="$service" -o wide
        
    else
        print_error "Deployment '$service' not found in namespace '$namespace'"
        return 1
    fi
    
    # Get pod info
    echo ""
    print_debug "Pod Status:"
    local pods=$(kubectl get pods -n "$namespace" -l app="$service" --no-headers -o custom-columns="NAME:.metadata.name" 2>/dev/null)
    
    if [ -n "$pods" ]; then
        kubectl get pods -n "$namespace" -l app="$service" -o wide
        
        # Check each pod
        echo "$pods" | while read pod; do
            echo ""
            print_debug "Pod Details: $pod"
            
            # Pod status
            local phase=$(kubectl get pod "$pod" -n "$namespace" -o jsonpath='{.status.phase}')
            local ready=$(kubectl get pod "$pod" -n "$namespace" -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}')
            echo "  Phase: $phase"
            echo "  Ready: $ready"
            
            # Check for restarts
            local restarts=$(kubectl get pod "$pod" -n "$namespace" -o jsonpath='{.status.containerStatuses[0].restartCount}')
            if [ "$restarts" -gt 0 ]; then
                print_warning "  Restart count: $restarts"
            fi
            
            # Show recent events for this pod
            print_debug "Recent Events for $pod:"
            kubectl get events -n "$namespace" --field-selector involvedObject.name="$pod" --sort-by='.lastTimestamp' | tail -n 5
            
            # If pod is not running, show more details
            if [ "$phase" != "Running" ] || [ "$ready" != "True" ]; then
                print_debug "Pod Description (last 20 lines):"
                kubectl describe pod "$pod" -n "$namespace" | tail -n 20
                
                print_debug "Recent Logs:"
                kubectl logs "$pod" -n "$namespace" --tail=20 2>/dev/null || print_warning "Could not retrieve logs"
            fi
            
            if [ "$DETAILED" = true ]; then
                print_debug "Full Pod Description:"
                kubectl describe pod "$pod" -n "$namespace"
            fi
        done
    else
        print_error "No pods found for service '$service'"
    fi
    
    # Check service
    echo ""
    print_debug "Service Configuration:"
    if kubectl get service "$service" -n "$namespace" >/dev/null 2>&1; then
        kubectl get service "$service" -n "$namespace" -o wide
        
        if [ "$DETAILED" = true ]; then
            echo ""
            print_debug "Service Description:"
            kubectl describe service "$service" -n "$namespace"
        fi
    else
        print_warning "Service '$service' not found"
    fi
    
    # Check configmaps
    echo ""
    print_debug "ConfigMaps:"
    kubectl get configmaps -n "$namespace" -l app="$service" 2>/dev/null || print_info "No configmaps found for $service"
    
    # Check secrets
    echo ""
    print_debug "Secrets:"
    kubectl get secrets -n "$namespace" -l app="$service" 2>/dev/null || print_info "No specific secrets found for $service"
}

# Function to debug all services
debug_all_services() {
    local namespace=$1
    
    print_header "Debugging All Services"
    
    # Get all deployments
    local deployments=$(kubectl get deployments -n "$namespace" --no-headers -o custom-columns="NAME:.metadata.name" 2>/dev/null)
    
    if [ -z "$deployments" ]; then
        print_error "No deployments found in namespace '$namespace'"
        return 1
    fi
    
    print_info "Found deployments: $(echo $deployments | tr '\n' ' ')"
    
    # Quick overview
    echo ""
    print_debug "Quick Overview:"
    kubectl get deployments,pods,services -n "$namespace" -o wide
    
    echo ""
    print_debug "Resource Usage (if available):"
    kubectl top pods -n "$namespace" 2>/dev/null || print_info "Resource metrics not available"
    
    # Check problematic pods
    echo ""
    print_debug "Problematic Pods:"
    local problem_pods=$(kubectl get pods -n "$namespace" --field-selector=status.phase!=Running --no-headers 2>/dev/null | wc -l)
    if [ "$problem_pods" -gt 0 ]; then
        kubectl get pods -n "$namespace" --field-selector=status.phase!=Running -o wide
        
        # Show details for problematic pods
        kubectl get pods -n "$namespace" --field-selector=status.phase!=Running --no-headers -o custom-columns="NAME:.metadata.name" | while read pod; do
            echo ""
            print_debug "Problem Pod: $pod"
            kubectl describe pod "$pod" -n "$namespace" | tail -n 10
        done
    else
        print_success "No problematic pods found"
    fi
    
    # Recent events
    echo ""
    print_debug "Recent Events (last 10):"
    kubectl get events -n "$namespace" --sort-by='.lastTimestamp' | tail -n 10
    
    # If detailed mode, debug each service
    if [ "$DETAILED" = true ]; then
        echo "$deployments" | while read deployment; do
            echo ""
            debug_service "$deployment" "$namespace"
        done
    fi
}

# Function to check cluster-wide issues
debug_cluster() {
    print_header "Cluster-wide Debug Information"
    
    print_debug "Cluster Info:"
    kubectl cluster-info
    
    echo ""
    print_debug "Node Status:"
    kubectl get nodes -o wide
    
    echo ""
    print_debug "Namespace Status:"
    kubectl get namespaces
    
    echo ""
    print_debug "System Pods:"
    kubectl get pods -n kube-system | grep -E "(coredns|traefik|local-path)" || print_info "System pods look normal"
    
    if [ "$DETAILED" = true ]; then
        echo ""
        print_debug "Cluster Events (last 20):"
        kubectl get events --all-namespaces --sort-by='.lastTimestamp' | tail -n 20
    fi
}

# Main execution
debug_cluster

echo ""
if [ -n "$SERVICE" ]; then
    debug_service "$SERVICE" "$NAMESPACE"
else
    debug_all_services "$NAMESPACE"
fi

echo ""
print_success "Debug information collection completed"
print_info "If issues persist, check:"
print_info "  • Resource limits and requests"
print_info "  • Image pull policies and registry access"
print_info "  • Network policies and service mesh configuration"
print_info "  • Secret and ConfigMap configurations"
