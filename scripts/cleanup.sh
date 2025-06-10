#!/bin/bash

# Cleanup script for Naytife development environment
set -e

echo "🧹 Cleaning up Naytife Development Environment..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}This will delete the k3d cluster and all data. Are you sure? (y/N)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    echo "🗑️  Deleting k3d cluster..."
    k3d cluster delete naytife 2>/dev/null || echo "Cluster already deleted or doesn't exist"
    
    echo "🗑️  Removing Docker images..."
    docker rmi naytife-registry.localhost:5000/naytife/auth-handler:latest 2>/dev/null || true
    docker rmi naytife-registry.localhost:5000/naytife/backend:latest 2>/dev/null || true
    docker rmi naytife-registry.localhost:5000/naytife/cloud-build:latest 2>/dev/null || true
    
    echo "🗑️  Stopping Tilt (if running)..."
    pkill -f "tilt up" 2>/dev/null || true
    
    echo -e "${GREEN}✅ Cleanup complete!${NC}"
    echo ""
    echo "To redeploy:"
    echo "  ./scripts/deploy-all-services.sh"
else
    echo "Cleanup cancelled."
fi
