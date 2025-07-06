#!/bin/bash

# Performance Benchmarking Tool for Naytife Platform
# Usage: ./benchmark.sh [environment] [--service=SERVICE] [--duration=30s] [--concurrency=10]

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
print_header() { echo -e "${CYAN}⚡ $1${NC}"; }
print_metric() { echo -e "${MAGENTA}📊 $1${NC}"; }

# Default values
ENVIRONMENT="local"
SPECIFIC_SERVICE=""
DURATION="30s"
CONCURRENCY=10
OUTPUT_DIR="/tmp/naytife-benchmarks"

show_usage() {
    echo -e "${CYAN}⚡ Naytife Performance Benchmark Tool${NC}"
    echo "===================================="
    echo "Usage: $0 [environment] [options]"
    echo ""
    echo -e "${YELLOW}Environments:${NC} local, staging, production"
    echo ""
    echo -e "${YELLOW}Options:${NC}"
    echo "  --service=SERVICE      Benchmark specific service"
    echo "  --duration=DURATION    Test duration (default: 30s)"
    echo "  --concurrency=N        Concurrent requests (default: 10)"
    echo "  --output-dir=DIR       Output directory for results"
    echo "  --help, -h            Show this help"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 local                                    # Benchmark all services"
    echo "  $0 staging --service=backend --duration=60s # Benchmark backend for 60s"
    echo "  $0 local --concurrency=20                   # Use 20 concurrent requests"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --service=*)
            SPECIFIC_SERVICE="${1#*=}"
            shift
            ;;
        --duration=*)
            DURATION="${1#*=}"
            shift
            ;;
        --concurrency=*)
            CONCURRENCY="${1#*=}"
            shift
            ;;
        --output-dir=*)
            OUTPUT_DIR="${1#*=}"
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

print_header "Performance Benchmarks for $ENVIRONMENT Environment"
echo "===================================================="

# Check prerequisites
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    print_error "curl is not installed"
    exit 1
fi

# Check for hey (HTTP load testing tool)
if ! command -v hey &> /dev/null; then
    print_warning "hey is not installed. Installing via go..."
    if command -v go &> /dev/null; then
        go install github.com/rakyll/hey@latest
        export PATH=$PATH:$(go env GOPATH)/bin
    else
        print_error "Neither hey nor go is installed. Please install hey or go first."
        print_info "Install hey: https://github.com/rakyll/hey"
        exit 1
    fi
fi

if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_error "Namespace '$NAMESPACE' does not exist"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
REPORT_FILE="$OUTPUT_DIR/benchmark_${ENVIRONMENT}_${TIMESTAMP}.md"

print_success "Prerequisites check passed"
print_info "Results will be saved to: $REPORT_FILE"

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

# Function to initialize report
init_report() {
    cat > "$REPORT_FILE" << EOF
# Naytife Performance Benchmark Report

**Environment:** $ENVIRONMENT  
**Date:** $(date)  
**Duration:** $DURATION  
**Concurrency:** $CONCURRENCY  
**Namespace:** $NAMESPACE  

## Test Configuration

- Load testing tool: hey
- Test duration: $DURATION
- Concurrent requests: $CONCURRENCY
- Environment: $ENVIRONMENT

## Results

EOF
}

# Function to add section to report
add_to_report() {
    echo "$1" >> "$REPORT_FILE"
}

# Function to benchmark a service
benchmark_service() {
    local service=$1
    local namespace=$2
    local port=$3
    local endpoint="${4:-/health}"
    
    print_header "Benchmarking $service"
    
    local service_url=$(get_service_url "$service" "$namespace" "$port")
    local full_url="${service_url}${endpoint}"
    
    print_info "Testing URL: $full_url"
    
    # Check if service is accessible
    if ! curl -s --max-time 5 "$full_url" >/dev/null 2>&1; then
        print_warning "Service $service is not accessible, skipping benchmark"
        add_to_report "### $service"
        add_to_report ""
        add_to_report "⚠️ **Service not accessible** - Skipped"
        add_to_report ""
        return 1
    fi
    
    print_info "Starting benchmark for $service..."
    
    # Run hey benchmark
    local output_file="$OUTPUT_DIR/${service}_${TIMESTAMP}.txt"
    
    # Use appropriate flags based on environment
    local extra_flags=""
    if [ "$ENVIRONMENT" != "local" ]; then
        extra_flags="-disable-keepalive"
    fi
    
    if hey -z "$DURATION" -c "$CONCURRENCY" $extra_flags "$full_url" > "$output_file" 2>&1; then
        print_success "Benchmark completed for $service"
        
        # Parse results
        local total_requests=$(grep "Total:" "$output_file" | awk '{print $2}')
        local success_rate=$(grep "Success rate:" "$output_file" | awk '{print $3}' || echo "N/A")
        local avg_response=$(grep "Average:" "$output_file" | awk '{print $2}')
        local fastest=$(grep "Fastest:" "$output_file" | awk '{print $2}')
        local slowest=$(grep "Slowest:" "$output_file" | awk '{print $2}')
        local rps=$(grep "Requests/sec:" "$output_file" | awk '{print $2}')
        
        # Display results
        print_metric "Service: $service"
        print_metric "Total requests: $total_requests"
        print_metric "Success rate: $success_rate"
        print_metric "Average response time: $avg_response"
        print_metric "Fastest response: $fastest"
        print_metric "Slowest response: $slowest"
        print_metric "Requests per second: $rps"
        
        # Add to report
        add_to_report "### $service"
        add_to_report ""
        add_to_report "**URL:** \`$full_url\`"
        add_to_report ""
        add_to_report "| Metric | Value |"
        add_to_report "|--------|-------|"
        add_to_report "| Total Requests | $total_requests |"
        add_to_report "| Success Rate | $success_rate |"
        add_to_report "| Average Response Time | $avg_response |"
        add_to_report "| Fastest Response | $fastest |"
        add_to_report "| Slowest Response | $slowest |"
        add_to_report "| Requests per Second | $rps |"
        add_to_report ""
        add_to_report "<details>"
        add_to_report "<summary>Full Results</summary>"
        add_to_report ""
        add_to_report "\`\`\`"
        cat "$output_file" >> "$REPORT_FILE"
        add_to_report "\`\`\`"
        add_to_report ""
        add_to_report "</details>"
        add_to_report ""
        
    else
        print_error "Benchmark failed for $service"
        add_to_report "### $service"
        add_to_report ""
        add_to_report "❌ **Benchmark failed**"
        add_to_report ""
    fi
    
    echo ""
}

# Function to get resource usage
get_resource_usage() {
    local namespace=$1
    
    print_header "Resource Usage During Benchmark"
    
    add_to_report "## Resource Usage"
    add_to_report ""
    
    # Check if metrics server is available
    if kubectl top nodes >/dev/null 2>&1; then
        print_info "Node resource usage:"
        kubectl top nodes
        
        print_info "Pod resource usage in $namespace:"
        kubectl top pods -n "$namespace"
        
        add_to_report "### Node Resource Usage"
        add_to_report ""
        add_to_report "\`\`\`"
        kubectl top nodes >> "$REPORT_FILE"
        add_to_report "\`\`\`"
        add_to_report ""
        
        add_to_report "### Pod Resource Usage"
        add_to_report ""
        add_to_report "\`\`\`"
        kubectl top pods -n "$namespace" >> "$REPORT_FILE"
        add_to_report "\`\`\`"
        add_to_report ""
        
    else
        print_warning "Metrics server not available - cannot show resource usage"
        add_to_report "⚠️ Metrics server not available"
        add_to_report ""
    fi
}

# Initialize report
init_report

# Get initial resource usage
get_resource_usage "$NAMESPACE"

# Define service benchmarks
# Use typeset which works in both bash and zsh
typeset -A SERVICES

# Populate service mappings
SERVICES[backend]="8000:/health"
SERVICES[auth-handler]="8080:/health"
SERVICES[hydra]="4444:/health/ready"
SERVICES[oathkeeper]="4456:/health/ready"

# Add optional services if they exist
if kubectl get deployment store-deployer -n "$NAMESPACE" >/dev/null 2>&1; then
    SERVICES["store-deployer"]="8090:/health"
fi

if kubectl get deployment template-registry -n "$NAMESPACE" >/dev/null 2>&1; then
    SERVICES["template-registry"]="8091:/health"
fi

# Run benchmarks
if [ -n "$SPECIFIC_SERVICE" ]; then
    if [[ ! " ${!SERVICES[@]} " =~ " $SPECIFIC_SERVICE " ]]; then
        print_error "Unknown service: $SPECIFIC_SERVICE"
        exit 1
    fi
    
    IFS=':' read -r port endpoint <<< "${SERVICES[$SPECIFIC_SERVICE]}"
    benchmark_service "$SPECIFIC_SERVICE" "$NAMESPACE" "$port" "$endpoint"
else
    for service in "${!SERVICES[@]}"; do
        IFS=':' read -r port endpoint <<< "${SERVICES[$service]}"
        benchmark_service "$service" "$NAMESPACE" "$port" "$endpoint"
    done
fi

# Get final resource usage
add_to_report "## Post-Benchmark Resource Usage"
add_to_report ""

if kubectl top nodes >/dev/null 2>&1; then
    add_to_report "### Node Resource Usage (After)"
    add_to_report ""
    add_to_report "\`\`\`"
    kubectl top nodes >> "$REPORT_FILE"
    add_to_report "\`\`\`"
    add_to_report ""
    
    add_to_report "### Pod Resource Usage (After)"
    add_to_report ""
    add_to_report "\`\`\`"
    kubectl top pods -n "$NAMESPACE" >> "$REPORT_FILE"
    add_to_report "\`\`\`"
    add_to_report ""
fi

# Add recommendations
add_to_report "## Recommendations"
add_to_report ""
add_to_report "- Monitor response times under normal load"
add_to_report "- Set up alerting for response time thresholds"
add_to_report "- Consider horizontal pod autoscaling if response times degrade"
add_to_report "- Review resource requests and limits based on usage patterns"
add_to_report ""

print_success "Benchmark completed!"
print_info "Full report available at: $REPORT_FILE"

# Show quick summary
echo ""
print_header "Quick Summary"
if [ -n "$SPECIFIC_SERVICE" ]; then
    grep -A 10 "### $SPECIFIC_SERVICE" "$REPORT_FILE" | grep "|" | head -7
else
    for service in "${!SERVICES[@]}"; do
        echo ""
        print_info "$service:"
        grep -A 10 "### $service" "$REPORT_FILE" | grep "| Requests per Second" | head -1
    done
fi
