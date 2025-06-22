#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📦 Loading Docker Images into k3d Cluster${NC}"
echo "=============================================="

CLUSTER_NAME=${1:-naytife}

# Check if cluster exists
if ! k3d cluster list | grep -q "$CLUSTER_NAME"; then
    echo -e "${RED}❌ k3d cluster '$CLUSTER_NAME' not found${NC}"
    echo "Create it first with: ./scripts/create-cluster.sh"
    exit 1
fi

echo -e "${BLUE}Loading images into cluster: $CLUSTER_NAME${NC}"

# Images to load
IMAGES=(
    "naytife/backend:latest"
    "naytife/auth-handler:latest"
    "template-registry:latest"
    "store-deployer:latest"
)

for image in "${IMAGES[@]}"; do
    echo -e "\n${YELLOW}📥 Loading $image...${NC}"
    
    # Check if image exists locally
    if ! docker image inspect "$image" >/dev/null 2>&1; then
        echo -e "${RED}❌ Image $image not found locally${NC}"
        echo "Build it first with: ./scripts/build-images.sh"
        continue
    fi
    
    # Load image into k3d cluster
    k3d image import "$image" -c "$CLUSTER_NAME"
    echo -e "${GREEN}✅ Loaded $image${NC}"
done

echo -e "\n${GREEN}🎉 All images loaded successfully!${NC}"

# Verify images in cluster
echo -e "\n${BLUE}📊 Verifying images in cluster...${NC}"
echo -e "${YELLOW}Note: Images will be available to pods once deployed${NC}"

echo -e "\n${BLUE}📝 Next Steps:${NC}"
echo "  1. Deploy services: ./scripts/deploy-all.sh"
echo "  2. Check status: ./scripts/status.sh"
