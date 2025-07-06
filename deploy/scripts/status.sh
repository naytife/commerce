#!/bin/bash

# Enhanced Status Script for Kustomize-based deployments
# Usage: ./status.sh [environment] [--detailed]

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
print_header() { echo -e "${CYAN}📊 $1${NC}"; }

# Default values
ENVIRONMENT="local"
DETAILED=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --detailed)
            DETAILED=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [environment] [--detailed]"
            echo "Environments: local, staging, production"
            echo "Options:"
            echo "  --detailed    Show detailed resource information"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

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
esac

print_header "Naytife Platform Status - $ENVIRONMENT Environment"
echo "============================================="

# Check if cluster is accessible
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

print_success "Connected to Kubernetes cluster"

# Check if namespace exists
if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_error "Namespace '$NAMESPACE' not found"
    exit 1
fi

print_success "Namespace '$NAMESPACE' exists"

# Check auth namespace
if kubectl get namespace "$AUTH_NAMESPACE" >/dev/null 2>&1; then
    AUTH_AVAILABLE=true
    print_success "Auth namespace '$AUTH_NAMESPACE' exists"
else
    AUTH_AVAILABLE=false
    print_warning "Auth namespace '$AUTH_NAMESPACE' not found - auth services won't be checked"
fi

echo ""
print_header "Cluster Information"
kubectl cluster-info | head -n 2

echo ""
print_header "Deployment Status"
kubectl get deployments -n "$NAMESPACE" -o wide

echo ""
print_header "Pod Status"
kubectl get pods -n "$NAMESPACE" -o wide

echo ""
print_header "Service Status"
kubectl get services -n "$NAMESPACE" -o wide

# Show auth services if available
if [ "$AUTH_AVAILABLE" = true ]; then
    echo ""
    print_header "Auth Services Status (in $AUTH_NAMESPACE)"
    echo ""
    print_info "Deployments:"
    kubectl get deployments -n "$AUTH_NAMESPACE" -o wide
    echo ""
    print_info "Pods:"
    kubectl get pods -n "$AUTH_NAMESPACE" -o wide
    echo ""
    print_info "Services:"
    kubectl get services -n "$AUTH_NAMESPACE" -o wide
fi

if [ "$DETAILED" = true ]; then
    echo ""
    print_header "ConfigMaps"
    kubectl get configmaps -n "$NAMESPACE"
    
    echo ""
    print_header "Secrets"
    kubectl get secrets -n "$NAMESPACE"
    
    echo ""
    print_header "Persistent Volumes"
    kubectl get pv,pvc -n "$NAMESPACE"
    
    echo ""
    print_header "Ingress"
    kubectl get ingress -n "$NAMESPACE" 2>/dev/null || print_info "No ingress resources found"
fi

echo ""
print_header "Service Health Checks"

# Function to check pod readiness
check_pod_health() {
    local pod_name=$1
    local namespace=$2
    
    # Get pod status
    local ready=$(kubectl get pod "$pod_name" -n "$namespace" -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' 2>/dev/null || echo "Unknown")
    local phase=$(kubectl get pod "$pod_name" -n "$namespace" -o jsonpath='{.status.phase}' 2>/dev/null || echo "Unknown")
    
    if [ "$ready" = "True" ] && [ "$phase" = "Running" ]; then
        echo -e "${GREEN}✅ $pod_name${NC}"
    elif [ "$phase" = "Pending" ]; then
        echo -e "${YELLOW}⏳ $pod_name (Pending)${NC}"
    else
        echo -e "${RED}❌ $pod_name (Phase: $phase, Ready: $ready)${NC}"
    fi
}

# Check all pods in the namespace
echo "Pod Health:"
kubectl get pods -n "$NAMESPACE" --no-headers -o custom-columns="NAME:.metadata.name" | while read pod; do
    check_pod_health "$pod" "$NAMESPACE"
done

echo ""
print_header "Resource Usage"

# Show resource usage if metrics server is available
if kubectl top nodes >/dev/null 2>&1; then
    echo "Node Resource Usage:"
    kubectl top nodes
    echo ""
    echo "Pod Resource Usage (in $NAMESPACE):"
    kubectl top pods -n "$NAMESPACE" 2>/dev/null || print_info "Pod metrics not available"
else
    print_warning "Metrics server not available - cannot show resource usage"
fi

echo ""
print_header "Recent Events"
kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -n 10

# Environment-specific checks
case $ENVIRONMENT in
    local)
        echo ""
        print_header "Local Development Information"
        print_info "NodePort Services:"
        kubectl get services -n "$NAMESPACE" -o wide | grep NodePort || print_info "No NodePort services found"
        ;;
    staging|production)
        echo ""
        print_header "External Access"
        kubectl get ingress -n "$NAMESPACE" -o wide 2>/dev/null || print_info "No ingress resources found"
        ;;
esac

echo ""
print_success "Status check completed for $ENVIRONMENT environment"
