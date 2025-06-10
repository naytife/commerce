#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Building Docker Images for Naytife Services${NC}"
echo "=================================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Function to build and tag image
build_image() {
    local service_name=$1
    local build_context=$2
    local dockerfile_path=$3
    
    echo -e "\n${BLUE}Building $service_name...${NC}"
    
    if [ ! -d "$build_context" ]; then
        echo -e "${RED}❌ Build context directory not found: $build_context${NC}"
        return 1
    fi
    
    if [ ! -f "$dockerfile_path" ]; then
        echo -e "${RED}❌ Dockerfile not found: $dockerfile_path${NC}"
        return 1
    fi
    
    docker build \
        -t "naytife/$service_name:latest" \
        -t "naytife/$service_name:dev" \
        -f "$dockerfile_path" \
        "$build_context"
    
    echo -e "${GREEN}✅ Built naytife/$service_name:latest${NC}"
}

# Build backend service
echo -e "\n${YELLOW}📦 Building Backend API${NC}"
build_image "backend" "$PROJECT_ROOT/backend" "$PROJECT_ROOT/backend/Dockerfile"

# Build auth handler
echo -e "\n${YELLOW}🔐 Building Auth Handler${NC}"
build_image "auth-handler" "$PROJECT_ROOT/auth/authentication-handler" "$PROJECT_ROOT/auth/authentication-handler/Dockerfile"

# Build cloud build service
echo -e "\n${YELLOW}🏗️  Building Cloud Build Service${NC}"
build_image "cloud-build" "$PROJECT_ROOT/cloud-build" "$PROJECT_ROOT/cloud-build/Dockerfile"

echo -e "\n${GREEN}🎉 All images built successfully!${NC}"
echo -e "\n${BLUE}📊 Built Images:${NC}"
docker images | grep "naytife/" | head -6

echo -e "\n${BLUE}📝 Next Steps:${NC}"
echo "  1. Load images into k3d: ./scripts/load-images.sh"
echo "  2. Deploy services: ./scripts/deploy-all.sh"
echo "  3. Check status: ./scripts/status.sh"
