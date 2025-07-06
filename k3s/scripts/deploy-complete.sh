#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Naytife Complete Deployment Script${NC}"
echo "====================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Check prerequisites
echo -e "\n${BLUE}🔍 Step 1: Checking Prerequisites${NC}"
echo "--------------------------------"

MISSING_DEPS=()
for cmd in k3d kubectl docker; do
    if command -v "$cmd" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $cmd is installed${NC}"
    else
        echo -e "${RED}❌ $cmd is not installed${NC}"
        MISSING_DEPS+=($cmd)
    fi
done

if [ ${#MISSING_DEPS[@]} -ne 0 ]; then
    echo -e "\n${RED}Missing dependencies: ${MISSING_DEPS[*]}${NC}"
    echo "Please install missing dependencies and run this script again."
    exit 1
fi

# Check if .env file exists
if [ ! -f "$PROJECT_ROOT/.env" ]; then
    echo -e "\n${YELLOW}⚠️  No .env file found${NC}"
    if [ -f "$PROJECT_ROOT/.env.example" ]; then
        echo "Copying .env.example to .env..."
        cp "$PROJECT_ROOT/.env.example" "$PROJECT_ROOT/.env"
        echo -e "${YELLOW}Please edit .env with your configuration before continuing${NC}"
        read -p "Press Enter to continue once you've configured .env..."
    else
        echo -e "${RED}❌ No .env.example found either${NC}"
        exit 1
    fi
fi

# Step 2: Create cluster
echo -e "\n${BLUE}🏗️  Step 2: Creating k3d Cluster${NC}"
echo "-------------------------------"
"$SCRIPT_DIR/create-cluster.sh"

# Step 3: Build images
echo -e "\n${BLUE}🐳 Step 3: Building Docker Images${NC}"
echo "-------------------------------"
"$SCRIPT_DIR/build-images.sh"

# Step 4: Load images
echo -e "\n${BLUE}📦 Step 4: Loading Images into Cluster${NC}"
echo "------------------------------------"
"$SCRIPT_DIR/load-images.sh"

# Step 5: Deploy services
echo -e "\n${BLUE}🚀 Step 5: Deploying All Services${NC}"
echo "-------------------------------"
"$SCRIPT_DIR/deploy-all.sh"

# Step 6: Final status check
echo -e "\n${BLUE}📊 Step 7: Final Status Check${NC}"
echo "----------------------------"
sleep 5
"$SCRIPT_DIR/status.sh"

echo -e "\n${GREEN}🎉 Naytife deployment completed successfully!${NC}"
echo ""
echo -e "\n${BLUE}🌐 Your services are now available at:${NC}"
echo "============================================"
echo "  🔐 API Gateway:        http://127.0.0.1:8080"
echo "  🔙 Backend API:        http://127.0.0.1:8000"
echo "  🔑 Auth Handler:       http://127.0.0.1:3000"
echo "  🏗️  Template Registry:  http://127.0.0.1:9001"
echo "  🚀 Store Deployer:     http://127.0.0.1:9003"
echo "  🐘 PostgreSQL:         localhost:5432"
echo "  📊 Redis:              localhost:6379"
echo "  🆔 Hydra Public:       http://127.0.0.1:4444"
echo "  🆔 Hydra Admin:        http://127.0.0.1:4445"

echo -e "\n${BLUE}📚 Quick Links:${NC}"
echo "==============="
echo "  📖 API Documentation:   http://127.0.0.1:8080/v1/docs"
echo "  🎮 GraphQL Playground:  http://127.0.0.1:8080/playground"
echo "  ❤️  Health Checks:       http://127.0.0.1:8080/health"

echo -e "\n${BLUE}🛠️  Useful Commands:${NC}"
echo "==================="
echo "  • Check status:     ./scripts/status.sh"
echo "  • View logs:        ./scripts/logs.sh [service]"
echo "  • Cleanup:          ./scripts/cleanup.sh"
echo "  • Restart service:  kubectl rollout restart deployment/[service] -n [namespace]"

echo -e "\n${BLUE}🔧 Troubleshooting:${NC}"
echo "=================="
echo "  • If services are not responding, wait a few more minutes"
echo "  • Check pod status: kubectl get pods --all-namespaces"
echo "  • View specific logs: ./scripts/logs.sh [service-name]"
echo "  • Restart deployment: kubectl rollout restart deployment/[name] -n [namespace]"

echo -e "\n${GREEN}Happy coding! 🚀${NC}"
