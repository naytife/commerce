#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Creating k3d Cluster for Naytife${NC}"
echo "====================================="

CLUSTER_NAME=${1:-naytife}

# Check if cluster already exists
if k3d cluster list | grep -q "$CLUSTER_NAME"; then
    echo -e "${YELLOW}⚠️  Cluster '$CLUSTER_NAME' already exists${NC}"
    read -p "Do you want to delete and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}🗑️  Deleting existing cluster...${NC}"
        k3d cluster delete "$CLUSTER_NAME"
    else
        echo -e "${GREEN}✅ Using existing cluster${NC}"
        exit 0
    fi
fi

echo -e "\n${BLUE}🔧 Creating k3d cluster: $CLUSTER_NAME${NC}"

# Create k3d cluster with specific port mappings
k3d cluster create "$CLUSTER_NAME" \
    --port "8080:30080@server:0" \
    --port "5433:30432@server:0" \
    --port "6380:30379@server:0" \
    --port "3000:30300@server:0" \
    --port "8000:30800@server:0" \
    --port "9000:30900@server:0" \
    --port "4444:30444@server:0" \
    --port "4445:30445@server:0" \
    --agents 1 \
    --wait

echo -e "\n${GREEN}✅ Cluster created successfully!${NC}"

# Verify cluster
echo -e "\n${BLUE}📊 Cluster Information:${NC}"
k3d cluster list
kubectl cluster-info

echo -e "\n${BLUE}🔗 Port Mappings:${NC}"
echo "  🔐 Oathkeeper (API Gateway): http://127.0.0.1:8080"
echo "  🔙 Backend API:             http://127.0.0.1:8000"
echo "  🔐 Auth Handler:            http://127.0.0.1:3000"
echo "  🏗️  Cloud Build:             http://127.0.0.1:9000"
echo "  🐘 PostgreSQL:              localhost:5433"
echo "  📊 Redis:                   localhost:6380"
echo "  🆔 Hydra Public:            http://127.0.0.1:4444"
echo "  🆔 Hydra Admin:             http://127.0.0.1:4445"

echo -e "\n${BLUE}📝 Next Steps:${NC}"
echo "  1. Build images: ./scripts/build-images.sh"
echo "  2. Load images: ./scripts/load-images.sh"
echo "  3. Deploy services: ./scripts/deploy-all.sh"
