#!/bin/bash

# Test script to verify the new template system is working correctly
# This script will:
# 1. Check that both services are healthy
# 2. List available templates
# 3. Deploy a test store
# 4. Check deployment status

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing New Template System${NC}"
echo "==============================="

# Test variables
TEST_SHOP_ID="test-shop-$(date +%s)"
TEST_SUBDOMAIN="test-$(date +%s)"
TEMPLATE_NAME="template_1"
VERSION="1.0.0"

# Set up port forwarding
echo -e "\n${YELLOW}📡 Setting up port forwarding...${NC}"
kubectl port-forward -n naytife svc/template-registry 9001:9001 > /dev/null 2>&1 &
TEMPLATE_REGISTRY_PID=$!
kubectl port-forward -n naytife svc/store-deployer 9003:9003 > /dev/null 2>&1 &
STORE_DEPLOYER_PID=$!

# Wait for port forwarding to be ready
sleep 3

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"
    kill $TEMPLATE_REGISTRY_PID $STORE_DEPLOYER_PID 2>/dev/null || true
}
trap cleanup EXIT

# Test 1: Health checks
echo -e "\n${YELLOW}🏥 Testing health endpoints...${NC}"
TEMPLATE_REGISTRY_HEALTH=$(curl -s http://localhost:9001/health | jq -r '.status' 2>/dev/null || echo "failed")
STORE_DEPLOYER_HEALTH=$(curl -s http://localhost:9003/health | jq -r '.status' 2>/dev/null || echo "failed")

if [ "$TEMPLATE_REGISTRY_HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✅ Template Registry: healthy${NC}"
else
    echo -e "${RED}❌ Template Registry: $TEMPLATE_REGISTRY_HEALTH${NC}"
    exit 1
fi

if [ "$STORE_DEPLOYER_HEALTH" = "healthy" ]; then
    echo -e "${GREEN}✅ Store Deployer: healthy${NC}"
else
    echo -e "${RED}❌ Store Deployer: $STORE_DEPLOYER_HEALTH${NC}"
    exit 1
fi

# Test 2: List templates
echo -e "\n${YELLOW}📋 Testing template listing...${NC}"
TEMPLATES=$(curl -s http://localhost:9003/templates)
echo "Available templates: $TEMPLATES"

TEMPLATE_COUNT=$(echo "$TEMPLATES" | jq -r '.count' 2>/dev/null || echo "0")
if [ "$TEMPLATE_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✅ Templates available: $TEMPLATE_COUNT${NC}"
else
    echo -e "${RED}❌ No templates found${NC}"
    exit 1
fi

# Test 3: Deploy a test store
echo -e "\n${YELLOW}🚀 Testing store deployment...${NC}"
DEPLOY_PAYLOAD=$(cat <<EOF
{
    "shop_id": "$TEST_SHOP_ID",
    "subdomain": "$TEST_SUBDOMAIN",
    "template_name": "$TEMPLATE_NAME",
    "version": "$VERSION"
}
EOF
)

echo "Deploying test store with payload:"
echo "$DEPLOY_PAYLOAD"

DEPLOY_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$DEPLOY_PAYLOAD" \
    http://localhost:9003/deploy)

echo "Deployment response: $DEPLOY_RESPONSE"

DEPLOY_STATUS=$(echo "$DEPLOY_RESPONSE" | jq -r '.status' 2>/dev/null || echo "failed")
if [ "$DEPLOY_STATUS" = "success" ]; then
    echo -e "${GREEN}✅ Store deployment successful${NC}"
    DEPLOYED_URL=$(echo "$DEPLOY_RESPONSE" | jq -r '.url' 2>/dev/null || echo "unknown")
    echo -e "${GREEN}   Store URL: $DEPLOYED_URL${NC}"
else
    echo -e "${RED}❌ Store deployment failed: $DEPLOY_STATUS${NC}"
    echo "Response: $DEPLOY_RESPONSE"
    exit 1
fi

# Test 4: Check deployment status
echo -e "\n${YELLOW}📊 Testing deployment status...${NC}"
STATUS_RESPONSE=$(curl -s http://localhost:9003/status/$TEST_SUBDOMAIN)
echo "Status response: $STATUS_RESPONSE"

echo -e "\n${GREEN}🎉 All tests passed! New template system is working correctly.${NC}"
echo -e "${GREEN}✨ Migration completed successfully!${NC}"
