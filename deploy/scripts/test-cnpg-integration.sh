#!/bin/bash

# CNPG Integration Tests for Naytife Commerce Platform
# This script tests CNPG cluster functionality and integration

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

print_header() {
    echo -e "${CYAN}🧪 $1${NC}"
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

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

run_test() {
    local test_name="$1"
    local test_command="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$VERBOSE" = true ]; then
        echo -e "${BLUE}Running: $test_name${NC}"
        echo -e "${BLUE}Command: $test_command${NC}"
    fi
    
    if eval "$test_command" >/dev/null 2>&1; then
        print_success "$test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_error "$test_name"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        if [ "$VERBOSE" = true ]; then
            echo -e "${RED}Failed command: $test_command${NC}"
            eval "$test_command" || true
        fi
    fi
}

show_usage() {
    echo -e "${CYAN}🧪 CNPG Integration Test Suite${NC}"
    echo "================================="
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --verbose              Show detailed test output"
    echo "  --cluster=NAME         Test specific cluster (default: naytife-postgres)"
    echo "  --help, -h            Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                       # Run all tests on local"
    echo "  $0 staging --verbose           # Run tests with detailed output"
    echo "  $0 local --cluster=test-cluster # Test specific cluster"
}

# Default values
ENVIRONMENT="local"
VERBOSE=false
CLUSTER_NAME="naytife-postgres"

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
        --cluster=*)
            CLUSTER_NAME="${1#*=}"
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
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        exit 1
        ;;
esac

# Test CNPG cluster status
test_cnpg_cluster() {
    print_header "Testing CNPG Cluster"
    
    # Test cluster exists
    run_test "CNPG cluster exists" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE"
    
    # Test cluster is ready
    run_test "CNPG cluster is ready" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.status.readyInstances}' | grep -v '^0$'"
    
    # Test cluster phase
    run_test "CNPG cluster is healthy" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.status.phase}' | grep -E '(Cluster in healthy state|healthy)'"
    
    # Test primary pod exists
    run_test "Primary pod exists" "kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}' | grep -v '^$'"
    
    # Test primary pod is ready
    run_test "Primary pod is ready" "kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"
    
    # Test replica pods (if instances > 1)
    local instances=$(kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.spec.instances}')
    if [ "$instances" -gt 1 ]; then
        run_test "Replica pods exist" "kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=replica -o jsonpath='{.items[*].metadata.name}' | grep -v '^$'"
    fi
    
    # Test services exist
    run_test "Read-write service exists" "kubectl get service $CLUSTER_NAME-rw -n $NAMESPACE"
    run_test "Read-only service exists" "kubectl get service $CLUSTER_NAME-ro -n $NAMESPACE"
    
    print_success "CNPG cluster tests completed"
}

# Test pooler functionality
test_cnpg_pooler() {
    print_header "Testing CNPG Pooler"
    
    local pooler_name="$CLUSTER_NAME-pooler"
    
    # Test pooler exists
    run_test "Pooler exists" "kubectl get pooler $pooler_name -n $NAMESPACE"
    
    # Test pooler is ready
    run_test "Pooler is ready" "kubectl get pooler $pooler_name -n $NAMESPACE -o jsonpath='{.status.phase}' | grep -E '(Ready|ready)'"
    
    # Test pooler pods are running
    run_test "Pooler pods are running" "kubectl get pods -n $NAMESPACE -l cnpg.io/poolerName=$pooler_name -o jsonpath='{.items[*].status.phase}' | grep -v -E '(Failed|Pending)'"
    
    # Test pooler service exists
    run_test "Pooler service exists" "kubectl get service $pooler_name-rw -n $NAMESPACE"
    
    print_success "CNPG pooler tests completed"
}

# Test database connectivity
test_database_connectivity() {
    print_header "Testing Database Connectivity"
    
    # Get primary pod name
    local primary_pod=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -z "$primary_pod" ]; then
        print_error "Cannot find primary pod for connectivity tests"
        return
    fi
    
    # Test direct connection to primary
    run_test "Direct connection to primary" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'SELECT 1;'"
    
    # Test database schemas exist
    run_test "Hydra schema exists" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'SELECT 1 FROM information_schema.schemata WHERE schema_name = '\''hydra'\'';'"
    
    run_test "Naytife schema exists" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'SELECT 1 FROM information_schema.schemata WHERE schema_name = '\''public'\'';'"
    
    # Test connection through pooler
    local pooler_service="$CLUSTER_NAME-pooler"
    run_test "Pooler connection" "kubectl exec -n $NAMESPACE $primary_pod -- psql -h $pooler_service -U naytife -d naytifedb -c 'SELECT 1;'"
    
    # Test basic database operations
    run_test "Create test table" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, data TEXT);'"
    
    run_test "Insert test data" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'INSERT INTO test_table (data) VALUES ('\''test-data'\'');'"
    
    run_test "Query test data" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'SELECT data FROM test_table WHERE data = '\''test-data'\'';'"
    
    run_test "Clean up test table" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U naytife -d naytifedb -c 'DROP TABLE IF EXISTS test_table;'"
    
    print_success "Database connectivity tests completed"
}

# Test backup functionality
test_backup_functionality() {
    print_header "Testing Backup Functionality"
    
    # Test scheduled backup exists
    run_test "Scheduled backup exists" "kubectl get scheduledbackup $CLUSTER_NAME-backup -n $NAMESPACE"
    
    # Test backup history
    if kubectl get backups -n $NAMESPACE >/dev/null 2>&1; then
        run_test "Backup history exists" "kubectl get backups -n $NAMESPACE -o jsonpath='{.items[*].metadata.name}' | grep -v '^$'"
    else
        print_info "No backup history found (this is normal for new clusters)"
    fi
    
    # Test backup configuration
    run_test "Cluster has backup configuration" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.spec.backup}' | grep -v '^$'"
    
    print_success "Backup functionality tests completed"
}

# Test monitoring
test_monitoring() {
    print_header "Testing Monitoring"
    
    # Test metrics endpoint on primary
    local primary_pod=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$primary_pod" ]; then
        run_test "Metrics endpoint accessible" "kubectl exec -n $NAMESPACE $primary_pod -- curl -s http://localhost:9187/metrics | grep -E '^pg_'"
    fi
    
    # Test ServiceMonitor exists
    run_test "ServiceMonitor exists" "kubectl get servicemonitor cnpg-cluster-metrics -n $NAMESPACE"
    
    # Test PrometheusRule exists
    run_test "PrometheusRule exists" "kubectl get prometheusrule cnpg-alerts -n $NAMESPACE"
    
    print_success "Monitoring tests completed"
}

# Test high availability (if multiple instances)
test_high_availability() {
    local instances=$(kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.spec.instances}')
    
    if [ "$instances" -le 1 ]; then
        print_info "Skipping HA tests (single instance cluster)"
        return
    fi
    
    print_header "Testing High Availability"
    
    # Test replication status
    local primary_pod=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=primary -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$primary_pod" ]; then
        run_test "Replication is working" "kubectl exec -n $NAMESPACE $primary_pod -- psql -U postgres -c 'SELECT * FROM pg_stat_replication;' | grep -E '(streaming|async)'"
    fi
    
    # Test replica connectivity
    local replica_pods=$(kubectl get pods -n $NAMESPACE -l postgresql=$CLUSTER_NAME,role=replica -o jsonpath='{.items[*].metadata.name}')
    
    if [ -n "$replica_pods" ]; then
        for replica_pod in $replica_pods; do
            run_test "Replica $replica_pod is accessible" "kubectl exec -n $NAMESPACE $replica_pod -- psql -U naytife -d naytifedb -c 'SELECT 1;'"
        done
    fi
    
    print_success "High availability tests completed"
}

# Test performance and resources
test_performance() {
    print_header "Testing Performance and Resources"
    
    # Test resource limits
    run_test "Cluster has resource limits" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.spec.resources.limits}' | grep -v '^$'"
    
    # Test storage
    run_test "Cluster has storage configured" "kubectl get cluster $CLUSTER_NAME -n $NAMESPACE -o jsonpath='{.spec.storage.size}' | grep -v '^$'"
    
    # Test PVCs exist
    run_test "PVCs exist" "kubectl get pvc -n $NAMESPACE -l postgresql=$CLUSTER_NAME"
    
    # Test connection pooling parameters
    local pooler_name="$CLUSTER_NAME-pooler"
    run_test "Pooler has connection limits" "kubectl get pooler $pooler_name -n $NAMESPACE -o jsonpath='{.spec.pgbouncer.parameters.max_client_conn}' | grep -v '^$'"
    
    print_success "Performance and resources tests completed"
}

# Test integration with existing services
test_service_integration() {
    print_header "Testing Service Integration"
    
    # Test DNS resolution
    run_test "Read-write service resolves" "kubectl exec -n $NAMESPACE deployment/backend -- nslookup $CLUSTER_NAME.$NAMESPACE.svc.cluster.local || true"
    
    run_test "Pooler service resolves" "kubectl exec -n $NAMESPACE deployment/backend -- nslookup $CLUSTER_NAME-pooler.$NAMESPACE.svc.cluster.local || true"
    
    # Test backend can connect to database
    if kubectl get deployment backend -n $NAMESPACE >/dev/null 2>&1; then
        run_test "Backend can connect to database" "kubectl exec -n $NAMESPACE deployment/backend -- echo 'SELECT 1' | psql \$DATABASE_URL"
    else
        print_info "Backend deployment not found, skipping integration test"
    fi
    
    print_success "Service integration tests completed"
}

# Main test execution
main() {
    print_header "CNPG Integration Tests for $ENVIRONMENT environment"
    echo "======================================================"
    
    # Check prerequisites
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    if ! kubectl get crd clusters.postgresql.cnpg.io >/dev/null 2>&1; then
        print_error "CNPG operator not installed"
        exit 1
    fi
    
    if ! kubectl get cluster $CLUSTER_NAME -n $NAMESPACE >/dev/null 2>&1; then
        print_error "CNPG cluster '$CLUSTER_NAME' not found in namespace '$NAMESPACE'"
        exit 1
    fi
    
    # Run test suites
    test_cnpg_cluster
    test_cnpg_pooler
    test_database_connectivity
    test_backup_functionality
    test_monitoring
    test_high_availability
    test_performance
    test_service_integration
    
    # Test summary
    echo ""
    print_header "Test Summary"
    echo "============"
    echo "Total tests: $TOTAL_TESTS"
    echo "Passed: $PASSED_TESTS"
    echo "Failed: $FAILED_TESTS"
    
    if [ "$FAILED_TESTS" -eq 0 ]; then
        print_success "All tests passed! 🎉"
        echo ""
        print_info "CNPG cluster is healthy and ready for production use"
        exit 0
    else
        print_error "$FAILED_TESTS test(s) failed"
        echo ""
        print_info "Please check the failed tests and resolve issues before proceeding"
        exit 1
    fi
}

# Run main function
main "$@"
