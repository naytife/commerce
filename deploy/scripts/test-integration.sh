#!/bin/bash

# Integration Test Suite for Naytife Platform
# Usage: ./test-integration.sh [environment] [--verbose] [--service=SERVICE]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}🧪 $1${NC}"; }

# Default values
ENVIRONMENT="local"
VERBOSE=false
SPECIFIC_SERVICE=""

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

show_usage() {
    echo -e "${CYAN}🧪 Naytife Integration Test Suite${NC}"
    echo "================================="
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --verbose              Show detailed test output"
    echo "  --service=SERVICE      Test only specific service"
    echo "  --help, -h            Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                       # Run all tests on local"
    echo "  $0 staging --verbose           # Run tests with detailed output"
    echo "  $0 local --service=backend     # Test only backend service"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --service=*)
            SPECIFIC_SERVICE="${1#*=}"
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

# Set namespaces based on environment
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        AUTH_NAMESPACE="naytife-auth"
        BUILD_NAMESPACE="naytife-build"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        AUTH_NAMESPACE="naytife-auth-staging"
        BUILD_NAMESPACE="naytife-build-staging"
        ;;
    production)
        NAMESPACE="naytife-production"
        AUTH_NAMESPACE="naytife-auth-production"
        BUILD_NAMESPACE="naytife-build-production"
        ;;
esac

print_header "Integration Tests for $ENVIRONMENT Environment"
echo "================================================="

# Check prerequisites
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    print_error "curl is not installed"
    exit 1
fi

if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_error "Namespace '$NAMESPACE' does not exist"
    exit 1
fi

if ! kubectl get namespace "$AUTH_NAMESPACE" >/dev/null 2>&1; then
    print_warning "Auth namespace '$AUTH_NAMESPACE' does not exist - skipping auth tests"
    AUTH_NAMESPACE=""
fi

print_success "Prerequisites check passed"

echo ""
print_info "Starting integration tests..."
print_info "Target namespaces:"
print_info "  - Core services: $NAMESPACE"
if [ -n "$AUTH_NAMESPACE" ]; then
    print_info "  - Auth services: $AUTH_NAMESPACE"
fi

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_result="${3:-0}"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$VERBOSE" = true ]; then
        print_info "Running: $test_name"
        echo "Command: $test_command"
    else
        echo -n "Testing: $test_name... "
    fi
    
    local output
    local result
    
    if [ "$VERBOSE" = true ]; then
        if output=$(eval "$test_command" 2>&1); then
            result=0
        else
            result=$?
        fi
        echo "Output: $output"
        echo "Exit code: $result"
    else
        if output=$(eval "$test_command" 2>&1); then
            result=0
        else
            result=$?
        fi
    fi
    
    if [ "$result" -eq "$expected_result" ]; then
        if [ "$VERBOSE" = false ]; then
            echo -e "${GREEN}✅ PASS${NC}"
        else
            print_success "TEST PASSED"
        fi
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        if [ "$VERBOSE" = false ]; then
            echo -e "${RED}❌ FAIL${NC}"
            echo "  Expected exit code: $expected_result, got: $result"
            echo "  Output: $output"
        else
            print_error "TEST FAILED"
            print_error "Expected exit code: $expected_result, got: $result"
        fi
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Function to get service URL
get_service_url() {
    local service=$1
    local namespace=$2
    local port=$3
    
    case $ENVIRONMENT in
        local)
            # For local, use NodePort
            local node_port=$(kubectl get service "$service" -n "$namespace" -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null || echo "")
            if [ -n "$node_port" ]; then
                echo "http://localhost:$node_port"
            else
                # Try port-forward
                echo "http://localhost:$port"
            fi
            ;;
        staging)
            echo "https://staging-${service}.naytife.com"
            ;;
        production)
            echo "https://${service}.naytife.com"
            ;;
    esac
}

# Function to test service connectivity
test_service_connectivity() {
    local service=$1
    local namespace=$2
    local port=$3
    local health_path="${4:-/health}"
    local prefix="${ENVIRONMENT}-"
    local prefixed_service="${prefix}${service}"
    
    print_header "Testing $service connectivity"
    
    # Test if service exists
    run_test "$service service exists" "kubectl get service '$prefixed_service' -n '$namespace'"
    
    # Test if pods are running
    run_test "$service pods are running" "kubectl get pods -n '$namespace' -l app='$service' --field-selector=status.phase=Running --no-headers | wc -l | grep -v '^0$'"
    
    # Test if pods are ready
    run_test "$service pods are ready" "kubectl get pods -n '$namespace' -l app='$service' -o jsonpath='{.items[*].status.conditions[?(@.type==\"Ready\")].status}' | grep -v False"
    
    # Test HTTP connectivity (if applicable)
    if [ "$service" != "postgres" ] && [ "$service" != "redis" ]; then
        local service_url=$(get_service_url "$prefixed_service" "$namespace" "$port")
        
        if [ "$ENVIRONMENT" = "local" ]; then
            # For local testing, we might need port-forwarding
            print_info "Testing HTTP connectivity to $service_url$health_path"
            
            # Try direct service access first
            if kubectl get service "$prefixed_service" -n "$namespace" -o jsonpath='{.spec.type}' | grep -q NodePort; then
                run_test "$service HTTP response" "curl -s --max-time 10 '$service_url$health_path' | grep -E '(ok|healthy|up|200)' || curl -s --max-time 10 '$service_url$health_path' -w '%{http_code}' | grep -E '^(200|201|204)$'"
            else
                print_warning "Service $service is not exposed via NodePort, skipping HTTP test"
            fi
        else
            # For staging/production, test external URLs
            run_test "$service HTTP response" "curl -s --max-time 10 '$service_url$health_path' | grep -E '(ok|healthy|up|200)' || curl -s --max-time 10 '$service_url$health_path' -w '%{http_code}' | grep -E '^(200|201|204)$'"
        fi
    fi
}

# Function to test CNPG cluster
test_cnpg_cluster() {
    local namespace=$1
    
    print_header "Testing CNPG Cluster"
    
    # Check if CNPG cluster exists
    if ! kubectl get cluster naytife-postgres -n "$namespace" >/dev/null 2>&1; then
        print_warning "CNPG cluster not found, skipping CNPG tests"
        return
    fi
    
    # Test cluster is ready
    run_test "CNPG cluster is ready" "kubectl get cluster naytife-postgres -n '$namespace' -o jsonpath='{.status.readyInstances}' | grep -v '^0$'"
    
    # Test cluster is healthy
    run_test "CNPG cluster is healthy" "kubectl get cluster naytife-postgres -n '$namespace' -o jsonpath='{.status.phase}' | grep -E '(Cluster in healthy state|healthy)'"
    
    # Test primary pod exists and is ready
    run_test "CNPG primary pod is ready" "kubectl get pods -n '$namespace' -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"
    
    # Test services exist
    run_test "CNPG read-write service exists" "kubectl get service naytife-postgres -n '$namespace'"
    run_test "CNPG read-only service exists" "kubectl get service naytife-postgres-ro -n '$namespace'"
    
    # Test pooler if it exists
    if kubectl get pooler naytife-postgres-pooler -n "$namespace" >/dev/null 2>&1; then
        run_test "CNPG pooler is ready" "kubectl get pooler naytife-postgres-pooler -n '$namespace' -o jsonpath='{.status.phase}' | grep -E '(Ready|ready)'"
        run_test "CNPG pooler service exists" "kubectl get service naytife-postgres-pooler -n '$namespace'"
    fi
    
    # Test database connectivity
    local primary_pod=$(kubectl get pods -n "$namespace" -l postgresql=naytife-postgres,role=primary -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
    if [ -n "$primary_pod" ]; then
        run_test "CNPG database is accessible" "kubectl exec -n '$namespace' '$primary_pod' -- psql -U naytife -d naytifedb -c 'SELECT 1;'"
        run_test "CNPG schemas exist" "kubectl exec -n '$namespace' '$primary_pod' -- psql -U naytife -d naytifedb -c 'SELECT schema_name FROM information_schema.schemata WHERE schema_name IN ('\''hydra'\'', '\''public'\'');' | grep -E '(hydra|public)'"
    fi
}

# Function to test database connectivity
test_database_connectivity() {
    local namespace=$1
    local prefix="${ENVIRONMENT}-"
    
    print_header "Testing Database Connectivity"
    
    # Check if CNPG cluster exists first
    if kubectl get cluster naytife-postgres -n "$namespace" >/dev/null 2>&1; then
        test_cnpg_cluster "$namespace"
    else
        # Fall back to traditional postgres testing
        print_info "CNPG cluster not found, testing traditional PostgreSQL deployment"
        
        # Test postgres
        run_test "Postgres pod is running" "kubectl get pods -n '$namespace' -l app=postgres --field-selector=status.phase=Running --no-headers | wc -l | grep -v '^0$'"
        
        # Test postgres connectivity from within cluster
        run_test "Postgres is accessible" "kubectl exec -n '$namespace' deployment/${prefix}postgres -- pg_isready -h localhost -p 5432"
    fi
    
    # Test redis
    run_test "Redis pod is running" "kubectl get pods -n '$namespace' -l app=redis --field-selector=status.phase=Running --no-headers | wc -l | grep -v '^0$'"
    
    # Test redis connectivity
    run_test "Redis is accessible" "kubectl exec -n '$namespace' deployment/${prefix}redis -- redis-cli -a naytife-redis-2024 ping | grep PONG"
}

# Function to test authentication flow
test_auth_flow() {
    local auth_namespace=$1
    
    if [ -z "$auth_namespace" ]; then
        print_warning "Auth namespace not available - skipping auth tests"
        return
    fi
    
    print_header "Testing Authentication Flow"
    
    # Test Hydra
    test_service_connectivity "hydra-admin" "$auth_namespace" "4445" "/health/ready"
    
    # Test Oathkeeper
    test_service_connectivity "oathkeeper-api" "$auth_namespace" "4456" "/health/ready"
    
    # Test Auth Handler
    test_service_connectivity "auth-handler" "$auth_namespace" "3000" "/health"
}

# Function to test backend services
test_backend_services() {
    local namespace=$1
    
    print_header "Testing Backend Services"
    
    # Test main backend
    test_service_connectivity "backend" "$namespace" "8000" "/health"
    
    # Test store deployer
    if kubectl get deployment store-deployer -n "$namespace" >/dev/null 2>&1; then
        test_service_connectivity "store-deployer" "$namespace" "8090" "/health"
    else
        print_warning "Store deployer not found, skipping tests"
    fi
    
    # Test template registry
    if kubectl get deployment template-registry -n "$namespace" >/dev/null 2>&1; then
        test_service_connectivity "template-registry" "$namespace" "8091" "/health"
    else
        print_warning "Template registry not found, skipping tests"
    fi
}

# Function to test end-to-end workflow
test_e2e_workflow() {
    local namespace=$1
    
    print_header "Testing End-to-End Workflow"
    
    print_warning "E2E workflow tests not yet implemented"
    print_info "Future tests should include:"
    print_info "  • User registration flow"
    print_info "  • Authentication and authorization"
    print_info "  • API endpoint functionality"
    print_info "  • Data persistence verification"
}

# Main test execution
echo ""
print_info "Starting integration tests..."
print_info "Target namespace: $NAMESPACE"

if [ -n "$SPECIFIC_SERVICE" ]; then
    print_info "Testing specific service: $SPECIFIC_SERVICE"
    case $SPECIFIC_SERVICE in
        postgres|redis)
            test_database_connectivity "$NAMESPACE"
            ;;
        hydra|oathkeeper|auth-handler)
            test_auth_flow "$AUTH_NAMESPACE"
            ;;
        backend|store-deployer|template-registry)
            test_backend_services "$NAMESPACE"
            ;;
        *)
            print_error "Unknown service: $SPECIFIC_SERVICE"
            exit 1
            ;;
    esac
else
    # Run all tests
    test_database_connectivity "$NAMESPACE"
    test_auth_flow "$AUTH_NAMESPACE"
    test_backend_services "$NAMESPACE"
    test_e2e_workflow "$NAMESPACE"
fi

# Test summary
echo ""
print_header "Test Summary"
echo "============"
echo "Total tests: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED_TESTS"

if [ "$FAILED_TESTS" -eq 0 ]; then
    print_success "All tests passed! 🎉"
    exit 0
else
    print_error "$FAILED_TESTS test(s) failed"
    exit 1
fi
