#!/bin/bash

# Build script for template system microservices
# Usage: ./build-template-services.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo -e "${BLUE}🔨 Building Template System Microservices${NC}"
echo -e "${BLUE}=======================================${NC}"

# Services to build (new simplified architecture)
SERVICES=("template-registry" "store-deployer")

for service in "${SERVICES[@]}"; do
    echo -e "\n${YELLOW}📦 Building $service...${NC}"
    
    cd "$PROJECT_ROOT/services/$service"
    
    # Download dependencies
    echo -e "${BLUE}Downloading dependencies...${NC}"
    go mod tidy
    
    # Build Go binary first to check for compilation errors
    echo -e "${BLUE}Compiling Go binary...${NC}"
    go build -o $service main.go
    
    # Build Docker image
    echo -e "${BLUE}Building Docker image...${NC}"
    docker build -t $service:latest .
    
    # Load into k3s if available
    if command -v k3s >/dev/null 2>&1; then
        echo -e "${BLUE}Loading image into k3s...${NC}"
        k3s ctr images import <(docker save $service:latest)
    fi
    
    echo -e "${GREEN}✅ $service built successfully${NC}"
done

echo -e "\n${GREEN}🎉 All template system services built successfully!${NC}"

# Show image sizes
echo -e "\n${BLUE}📊 Image sizes:${NC}"
for service in "${SERVICES[@]}"; do
    size=$(docker images $service:latest --format "table {{.Size}}" | tail -n 1)
    echo -e "   $service: $size"
done

echo -e "\n${YELLOW}🚀 Ready to deploy with: ./migrate-to-template-system.sh${NC}"
