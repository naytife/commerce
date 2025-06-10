#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧹 Naytife Cleanup Script${NC}"
echo "========================="

CLUSTER_NAME=${1:-naytife}
FORCE=${2:-false}

if [ "$FORCE" != "true" ] && [ "$FORCE" != "-f" ] && [ "$FORCE" != "--force" ]; then
    echo -e "${YELLOW}⚠️  This will delete:${NC}"
    echo "  • k3d cluster: $CLUSTER_NAME"
    echo "  • All deployed services"
    echo "  • All persistent data"
    echo ""
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}ℹ️  Cleanup cancelled${NC}"
        exit 0
    fi
fi

echo -e "\n${BLUE}🗑️  Starting cleanup...${NC}"

# Check if cluster exists
if k3d cluster list | grep -q "$CLUSTER_NAME"; then
    echo -e "\n${YELLOW}🏗️  Deleting k3d cluster: $CLUSTER_NAME${NC}"
    k3d cluster delete "$CLUSTER_NAME"
    echo -e "${GREEN}✅ Cluster deleted${NC}"
else
    echo -e "${YELLOW}⚠️  Cluster $CLUSTER_NAME not found${NC}"
fi

# Optional: Clean up Docker images
echo -e "\n${BLUE}🐳 Docker Image Cleanup${NC}"
read -p "Delete Naytife Docker images? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🗑️  Removing Naytife Docker images...${NC}"
    
    # Remove naytife images
    docker images | grep "naytife/" | awk '{print $1":"$2}' | xargs -r docker rmi 2>/dev/null || true
    
    # Remove dangling images
    docker image prune -f >/dev/null 2>&1 || true
    
    echo -e "${GREEN}✅ Docker images cleaned${NC}"
fi

# Optional: Clean up build artifacts
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo -e "\n${BLUE}🧹 Build Artifacts Cleanup${NC}"
read -p "Clean build artifacts? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🗑️  Cleaning build artifacts...${NC}"
    
    # Clean backend build artifacts
    if [ -d "$PROJECT_ROOT/backend/bin" ]; then
        rm -rf "$PROJECT_ROOT/backend/bin"
        echo -e "${GREEN}✅ Cleaned backend/bin${NC}"
    fi
    
    # Clean cloud-build artifacts
    if [ -d "$PROJECT_ROOT/cloud-build/built_sites" ]; then
        rm -rf "$PROJECT_ROOT/cloud-build/built_sites"
        echo -e "${GREEN}✅ Cleaned cloud-build/built_sites${NC}"
    fi
    
    # Clean any dump files
    find "$PROJECT_ROOT" -name "dump.rdb" -type f -delete 2>/dev/null || true
    find "$PROJECT_ROOT" -name "*.log" -path "*/tmp/*" -type f -delete 2>/dev/null || true
    
    echo -e "${GREEN}✅ Build artifacts cleaned${NC}"
fi

echo -e "\n${GREEN}🎉 Cleanup completed!${NC}"
echo ""
echo -e "${BLUE}📝 To restart the environment:${NC}"
echo "  1. Create cluster: ./scripts/create-cluster.sh"
echo "  2. Build images: ./scripts/build-images.sh"
echo "  3. Load images: ./scripts/load-images.sh"
echo "  4. Deploy services: ./scripts/deploy-all.sh"
