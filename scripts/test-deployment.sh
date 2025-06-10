#!/bin/bash

# Quick deployment test script
set -e

echo "🧪 Testing Naytife Deployment..."

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test if k3d cluster exists
if k3d cluster list | grep -q "naytife"; then
    echo -e "${GREEN}✅ k3d cluster exists${NC}"
else
    echo -e "${RED}❌ k3d cluster not found${NC}"
    echo "Run: ./k3s/scripts/deploy-all.sh to create the cluster"
    exit 1
fi

# Test if all pods are running
echo "Checking pod status..."
kubectl get pods --all-namespaces

echo ""
echo "Testing service endpoints..."

# Test Hydra health
echo -n "Hydra health: "
if curl -s http://localhost:4444/health/ready | grep -q "ok"; then
    echo -e "${GREEN}✅ Ready${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

# Test Auth Handler
echo -n "Auth Handler: "
if curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/health 2>/dev/null | grep -q "200"; then
    echo -e "${GREEN}✅ Ready${NC}"
else
    echo -e "${YELLOW}⚠️  Not ready or no health endpoint${NC}"
fi

# Test Backend API
echo -n "Backend API: "
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/v1/docs 2>/dev/null | grep -q "200"; then
    echo -e "${GREEN}✅ Ready${NC}"
else
    echo -e "${YELLOW}⚠️  Not ready${NC}"
fi

# Test OAuth2 flow
echo -n "OAuth2 Authorization: "
if curl -s -o /dev/null -w "%{http_code}" "http://127.0.0.1:8080/oauth2/auth?client_id=4b41cd38-43ed-4e3a-9a88-bd384af21732&response_type=code&scope=openid&redirect_uri=http://127.0.0.1:8080/auth/callback" 2>/dev/null | grep -q "302\|200"; then
    echo -e "${GREEN}✅ OAuth2 flow accessible${NC}"
else
    echo -e "${RED}❌ OAuth2 flow not accessible${NC}"
fi

echo ""
echo "🔧 Quick commands:"
echo "  kubectl get pods --all-namespaces"
echo "  kubectl logs -n auth deployment/hydra -f"
echo "  kubectl logs -n backend deployment/backend -f"
echo "  tilt up  # For development with hot reload"

echo ""
echo -e "${GREEN}Testing complete!${NC}"
