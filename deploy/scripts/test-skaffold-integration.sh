#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Naytife Platform - Skaffold Integration Test${NC}"
echo -e "${YELLOW}================================================${NC}"
echo ""

# Function to log status
log_step() {
    echo -e "${CYAN}📋 $1...${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
log_step "Checking prerequisites"

if ! command_exists kubectl; then
    log_error "kubectl not installed"
    exit 1
fi

if ! command_exists skaffold; then
    log_error "Skaffold not installed"
    echo -e "${YELLOW}   Install with: brew install skaffold${NC}"
    exit 1
fi

if ! command_exists docker; then
    log_error "Docker not installed"
    exit 1
fi

log_success "All required tools are installed"

# Check k3s cluster
log_step "Checking k3s cluster connectivity"

if ! kubectl get nodes &> /dev/null; then
    log_error "Cannot connect to k3s cluster"
    echo -e "${YELLOW}   Try running: ./k3s/scripts/create-cluster.sh${NC}"
    exit 1
fi

NODE_COUNT=$(kubectl get nodes --no-headers | wc -l | tr -d ' ')
log_success "Connected to k3s cluster with $NODE_COUNT node(s)"

# Check Docker daemon
log_step "Checking Docker daemon"

if ! docker ps &> /dev/null; then
    log_error "Docker daemon not running"
    exit 1
fi

log_success "Docker daemon is running"

# Validate Skaffold configuration
log_step "Validating Skaffold configuration"

if ! skaffold diagnose &> /dev/null; then
    log_error "Skaffold configuration is invalid"
    echo -e "${YELLOW}   Run: skaffold diagnose${NC}"
    exit 1
fi

log_success "Skaffold configuration is valid"

# Test Kustomize rendering
log_step "Testing Kustomize manifest rendering"

if ! skaffold render --profile=local &> /dev/null; then
    log_error "Kustomize rendering failed"
    echo -e "${YELLOW}   Run: skaffold render --profile=local${NC}"
    exit 1
fi

log_success "Kustomize manifests render successfully"

# Test Docker builds (dry run)
log_step "Testing Docker builds (dry run)"

if ! skaffold build --dry-run --profile=local &> /dev/null; then
    log_error "Docker build test failed"
    echo -e "${YELLOW}   Run: skaffold build --dry-run --profile=local${NC}"
    exit 1
fi

log_success "All Docker builds pass validation"

# Check deployment scripts
log_step "Checking deployment scripts"

SCRIPTS=(
    "deploy/scripts/deploy-skaffold.sh"
    "deploy/scripts/skaffold-utils.sh"
)

for script in "${SCRIPTS[@]}"; do
    if [ -f "$script" ] && [ -x "$script" ]; then
        log_success "$script is ready"
    else
        log_warning "$script is missing or not executable"
    fi
done

# Test namespace creation
log_step "Testing namespace management"

NAMESPACE="naytife-local"

if kubectl get namespace "$NAMESPACE" &> /dev/null; then
    log_success "Namespace $NAMESPACE already exists"
else
    log_warning "Namespace $NAMESPACE does not exist (will be created during deployment)"
fi

echo ""
echo -e "${GREEN}🎉 Skaffold Integration Test Complete!${NC}"
echo ""
echo -e "${BLUE}📋 Summary:${NC}"
echo -e "   ✅ Prerequisites installed and working"
echo -e "   ✅ k3s cluster accessible"
echo -e "   ✅ Skaffold configuration valid"
echo -e "   ✅ Kustomize manifests render correctly"
echo -e "   ✅ Docker builds pass validation"
echo ""
echo -e "${BLUE}🚀 Ready to start development!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "   1. Start development: ${GREEN}./deploy/scripts/deploy-skaffold.sh dev${NC}"
echo -e "   2. Or build and run: ${GREEN}./deploy/scripts/deploy-skaffold.sh run${NC}"
echo -e "   3. Monitor with: ${GREEN}./deploy/scripts/skaffold-utils.sh status${NC}"
echo -e "   4. View logs: ${GREEN}./deploy/scripts/skaffold-utils.sh logs${NC}"
echo ""
echo -e "${CYAN}📖 For more information, see: deploy/tools/README.md${NC}"
