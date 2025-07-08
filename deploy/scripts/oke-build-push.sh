#!/bin/bash

# Build and Push ARM64 Images to OCIR
# This script builds multi-architecture images and pushes them to Oracle Container Image Registry

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}📦 $1${NC}"; }

# Configuration - UPDATE THESE VALUES
REGION="${OCI_REGION:-us-ashburn-1}"  # Change to your region
TENANCY="${OCI_TENANCY:-}"            # Your tenancy namespace
OCIR_ENDPOINT="$REGION.ocir.io"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

print_header "Building and Pushing ARM64 Images to OCIR"

# Check if variables are set
if [ -z "$TENANCY" ]; then
    print_error "TENANCY not set. Please set OCI_TENANCY environment variable or update this script."
    print_info "To find your tenancy namespace: oci os ns get"
    exit 1
fi

print_info "OCIR Endpoint: $OCIR_ENDPOINT"
print_info "Tenancy: $TENANCY"

# Check if logged into OCIR
if ! docker info | grep -q "$OCIR_ENDPOINT"; then
    print_warning "Not logged into OCIR. Please login first:"
    echo "docker login $OCIR_ENDPOINT"
    echo "Username: $TENANCY/<your-username>"
    echo "Password: <your-auth-token>"
    exit 1
fi

# Create buildx builder if it doesn't exist
if ! docker buildx ls | grep -q "arm64-builder"; then
    print_info "Creating multi-architecture builder"
    docker buildx create --name arm64-builder --use
fi

# Use the multi-arch builder
docker buildx use arm64-builder

# Build and push backend
print_header "Building Backend (ARM64)"
cd "$PROJECT_ROOT/backend"

BACKEND_IMAGE="$OCIR_ENDPOINT/$TENANCY/naytife/backend:arm64-latest"
print_info "Building: $BACKEND_IMAGE"

docker buildx build \
    --platform linux/arm64 \
    --push \
    -t "$BACKEND_IMAGE" \
    --build-arg TARGETPLATFORM=linux/arm64 \
    .

print_success "Backend image built and pushed"

# Build and push auth-handler
print_header "Building Auth Handler (ARM64)"
cd "$PROJECT_ROOT/auth/authentication-handler"

AUTH_IMAGE="$OCIR_ENDPOINT/$TENANCY/naytife/auth-handler:arm64-latest"
print_info "Building: $AUTH_IMAGE"

docker buildx build \
    --platform linux/arm64 \
    --push \
    -t "$AUTH_IMAGE" \
    --build-arg TARGETPLATFORM=linux/arm64 \
    .

print_success "Auth handler image built and pushed"

# Build and push template-registry
print_header "Building Template Registry (ARM64)"
cd "$PROJECT_ROOT/services/template-registry"

TEMPLATE_IMAGE="$OCIR_ENDPOINT/$TENANCY/naytife/template-registry:arm64-latest"
print_info "Building: $TEMPLATE_IMAGE"

docker buildx build \
    --platform linux/arm64 \
    --push \
    -t "$TEMPLATE_IMAGE" \
    --build-arg TARGETPLATFORM=linux/arm64 \
    .

print_success "Template registry image built and pushed"

# Build and push store-deployer
print_header "Building Store Deployer (ARM64)"
cd "$PROJECT_ROOT/services/store-deployer"

DEPLOYER_IMAGE="$OCIR_ENDPOINT/$TENANCY/naytife/store-deployer:arm64-latest"
print_info "Building: $DEPLOYER_IMAGE"

docker buildx build \
    --platform linux/arm64 \
    --push \
    -t "$DEPLOYER_IMAGE" \
    --build-arg TARGETPLATFORM=linux/arm64 \
    .

print_success "Store deployer image built and pushed"

print_header "Build Summary"
print_success "All ARM64 images built and pushed to OCIR:"
echo "  - $BACKEND_IMAGE"
echo "  - $AUTH_IMAGE"
echo "  - $TEMPLATE_IMAGE"
echo "  - $DEPLOYER_IMAGE"

print_info "Next steps:"
echo "1. Update deploy/overlays/oke-free/kustomization.yaml with your region and tenancy"
echo "2. Deploy to OKE: ./oke-deploy.sh"

# Return to original directory
cd "$PROJECT_ROOT"
