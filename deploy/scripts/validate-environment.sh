#!/bin/bash

# Environment Validation Script
# Usage: ./validate-environment.sh [environment]

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
print_header() { echo -e "${CYAN}🔍 $1${NC}"; }

# Function to print test result
print_test_result() {
    if [ $1 -eq 0 ]; then
        print_success "$2"
        return 0
    else
        print_error "$2"
        return 1
    fi
}

# Default environment
ENVIRONMENT=${1:-local}

# Validate environment parameter
case $ENVIRONMENT in
    local|staging|production)
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        print_info "Valid environments: local, staging, production"
        exit 1
        ;;
esac

print_header "Environment Validation for $ENVIRONMENT"
echo "==========================================="

# Set environment-specific variables
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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
OVERLAY_DIR="$DEPLOY_DIR/overlays/$ENVIRONMENT"
SECRET_DIR="$DEPLOY_DIR/secrets/$ENVIRONMENT"

# Test counter
TOTAL_TESTS=0
PASSED_TESTS=0

run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -n "Testing: $test_name... "
    
    if eval "$test_command" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}❌ FAIL${NC}"
        return 1
    fi
}

echo ""
print_header "Prerequisites Check"

# Check required tools
run_test "kubectl is installed" "command -v kubectl"
run_test "kustomize is installed" "command -v kustomize"
run_test "sops is installed" "command -v sops"

# Check cluster connectivity
run_test "Kubernetes cluster is accessible" "kubectl cluster-info"

echo ""
print_header "Configuration Validation"

# Check directory structure
run_test "Overlay directory exists" "[ -d '$OVERLAY_DIR' ]"
run_test "Secret directory exists" "[ -d '$SECRET_DIR' ]"

# Check kustomization files
run_test "Kustomization.yaml exists" "[ -f '$OVERLAY_DIR/kustomization.yaml' ]"

# Test kustomize build
run_test "Kustomize configuration builds" "cd '$OVERLAY_DIR' && kustomize build . > /dev/null"

# Check secrets
if [ -d "$SECRET_DIR" ]; then
    secret_files=$(find "$SECRET_DIR" -name "*.yaml" 2>/dev/null || true)
    if [ -n "$secret_files" ]; then
        for secret_file in $secret_files; do
            filename=$(basename "$secret_file")
            run_test "Secret $filename can be decrypted" "sops -d '$secret_file' > /dev/null"
        done
    else
        print_warning "No secret files found in $SECRET_DIR"
    fi
fi

echo ""
print_header "Environment-Specific Validation"

case $ENVIRONMENT in
    local)
        # Local environment specific checks
        run_test "Local cluster context is available" "kubectl config get-contexts | grep -q k3d"
        
        # Check if namespace exists
        if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
            run_test "Namespace $NAMESPACE exists" "true"
            
            # Check if services are running
            if kubectl get deployments -n "$NAMESPACE" >/dev/null 2>&1; then
                deployment_count=$(kubectl get deployments -n "$NAMESPACE" --no-headers | wc -l)
                if [ "$deployment_count" -gt 0 ]; then
                    run_test "Deployments exist in namespace" "true"
                    
                    # Check deployment readiness
                    ready_deployments=$(kubectl get deployments -n "$NAMESPACE" -o jsonpath='{.items[?(@.status.readyReplicas==@.status.replicas)].metadata.name}' | wc -w)
                    total_deployments=$(kubectl get deployments -n "$NAMESPACE" --no-headers | wc -l)
                    
                    if [ "$ready_deployments" -eq "$total_deployments" ]; then
                        run_test "All deployments are ready" "true"
                    else
                        run_test "All deployments are ready ($ready_deployments/$total_deployments)" "false"
                    fi
                else
                    print_warning "No deployments found in namespace $NAMESPACE"
                fi
            fi
        else
            print_warning "Namespace $NAMESPACE does not exist - environment not deployed"
        fi
        ;;
        
    staging|production)
        # Remote environment specific checks
        run_test "Remote cluster context is available" "kubectl config current-context | grep -v k3d"
        
        # Check if ingress controller is available
        run_test "Ingress controller is running" "kubectl get pods -A -l app.kubernetes.io/name=ingress-nginx | grep -q Running"
        
        # Check SSL/TLS certificates if ingress exists
        if kubectl get ingress -n "$NAMESPACE" >/dev/null 2>&1; then
            ingress_count=$(kubectl get ingress -n "$NAMESPACE" --no-headers | wc -l)
            if [ "$ingress_count" -gt 0 ]; then
                run_test "Ingress resources exist" "true"
            fi
        fi
        ;;
esac

echo ""
print_header "Service Health Validation"

# If namespace exists, check service health
if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    # Check if all pods are running
    pod_count=$(kubectl get pods -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$pod_count" -gt 0 ]; then
        running_pods=$(kubectl get pods -n "$NAMESPACE" --field-selector=status.phase=Running --no-headers 2>/dev/null | wc -l)
        run_test "All pods are running ($running_pods/$pod_count)" "[ '$running_pods' -eq '$pod_count' ]"
        
        # Check pod readiness
        ready_pods=$(kubectl get pods -n "$NAMESPACE" -o jsonpath='{.items[?(@.status.conditions[?(@.type=="Ready")].status=="True")].metadata.name}' 2>/dev/null | wc -w)
        run_test "All pods are ready ($ready_pods/$pod_count)" "[ '$ready_pods' -eq '$pod_count' ]"
    fi
    
    # Check services
    service_count=$(kubectl get services -n "$NAMESPACE" --no-headers 2>/dev/null | wc -l)
    if [ "$service_count" -gt 0 ]; then
        run_test "Services are defined ($service_count services)" "true"
    fi
else
    print_warning "Namespace $NAMESPACE does not exist - skipping service health checks"
fi

echo ""
print_header "Validation Summary"

if [ "$PASSED_TESTS" -eq "$TOTAL_TESTS" ]; then
    print_success "All tests passed! ($PASSED_TESTS/$TOTAL_TESTS)"
    print_success "Environment $ENVIRONMENT is properly configured and ready"
    exit 0
else
    failed_tests=$((TOTAL_TESTS - PASSED_TESTS))
    print_error "$failed_tests tests failed ($PASSED_TESTS/$TOTAL_TESTS passed)"
    print_error "Environment $ENVIRONMENT has configuration issues"
    exit 1
fi
