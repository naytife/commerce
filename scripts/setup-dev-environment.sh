#!/bin/bash

# Complete development environment setup script
# This script deploys the entire k3s cluster and sets up OAuth2 clients

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${BLUE}🚀 Naytife Development Environment Setup${NC}"
echo "=============================================="
echo ""

# Step 1: Deploy k3s cluster
echo -e "${YELLOW}📦 Step 1: Deploying k3s cluster...${NC}"
if [ -f "$PROJECT_ROOT/k3s/scripts/deploy-all.sh" ]; then
    cd "$PROJECT_ROOT"
    ./k3s/scripts/deploy-all.sh
else
    echo -e "${RED}❌ k3s deployment script not found${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}⏳ Waiting for cluster to be fully ready...${NC}"
sleep 10

# Step 2: Create OAuth2 clients
echo -e "${YELLOW}🔑 Step 2: Creating OAuth2 clients...${NC}"
if [ -f "$SCRIPT_DIR/create-hydra-clients.sh" ]; then
    "$SCRIPT_DIR/create-hydra-clients.sh"
else
    echo -e "${RED}❌ Hydra clients script not found${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Development environment setup complete!${NC}"
echo ""
echo -e "${BLUE}📋 Available Services:${NC}"
echo "=================================================="
echo -e "${YELLOW}🌐 API Gateway:${NC}      http://127.0.0.1:8080"
echo -e "${YELLOW}📚 API Docs:${NC}        http://127.0.0.1:8080/v1/docs"
echo -e "${YELLOW}🔧 Backend:${NC}          http://127.0.0.1:8000"
echo -e "${YELLOW}🔐 Hydra Public:${NC}     http://127.0.0.1:4444"
echo -e "${YELLOW}🆔 Hydra Admin:${NC}      http://127.0.0.1:4445"
echo ""
echo -e "${BLUE}🖥️  Dashboard Setup:${NC}"
echo "=================================================="
echo "1. Navigate to the dashboard directory:"
echo "   cd dashboard"
echo ""
echo "2. Install dependencies:"
echo "   npm install"
echo ""
echo "3. Start the dashboard:"
echo "   npm run dev"
echo ""
echo "4. Open dashboard at:"
echo "   http://localhost:5173"
echo ""
echo -e "${GREEN}✅ Ready for development!${NC}"
