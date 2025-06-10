#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📊 Naytife Services Status${NC}"
echo "=========================="

# Check if cluster exists
if ! kubectl cluster-info >/dev/null 2>&1; then
    echo -e "${RED}❌ No k3s cluster found or kubectl not connected${NC}"
    exit 1
fi

echo -e "\n${BLUE}🏗️  Cluster Information:${NC}"
kubectl cluster-info --context k3d-naytife 2>/dev/null || kubectl cluster-info

echo -e "\n${BLUE}📦 Namespaces:${NC}"
kubectl get namespaces -l app.kubernetes.io/part-of=naytife-platform

echo -e "\n${BLUE}🚀 Pod Status:${NC}"
echo "======================================"
kubectl get pods --all-namespaces -l app.kubernetes.io/part-of=naytife-platform -o wide

echo -e "\n${BLUE}🔌 Services:${NC}"
echo "======================================"
kubectl get services --all-namespaces -l app.kubernetes.io/part-of=naytife-platform

echo -e "\n${BLUE}📡 Deployments:${NC}"
echo "======================================"
kubectl get deployments --all-namespaces -l app.kubernetes.io/part-of=naytife-platform

echo -e "\n${BLUE}🔗 Service Health Check:${NC}"
echo "======================================"

# Function to check service health
check_service() {
    local name=$1
    local url=$2
    local timeout=${3:-5}
    
    echo -n "$name: "
    if timeout $timeout curl -s "$url" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Healthy${NC}"
    else
        echo -e "${RED}❌ Unhealthy or not ready${NC}"
    fi
}

check_service "🔐 Oathkeeper  " "http://127.0.0.1:8080/health"
check_service "🔙 Backend     " "http://127.0.0.1:8000/health"
check_service "🔑 Auth Handler" "http://127.0.0.1:3000/health"
check_service "🏗️  Cloud Build " "http://127.0.0.1:9000/health"
check_service "🆔 Hydra Public" "http://127.0.0.1:4444/health/alive"

echo -e "\n${BLUE}🔗 Service URLs:${NC}"
echo "======================================"
echo "  🔐 API Gateway:    http://127.0.0.1:8080"
echo "  🔙 Backend API:    http://127.0.0.1:8000"
echo "  🔑 Auth Handler:   http://127.0.0.1:3000"
echo "  🏗️  Cloud Build:    http://127.0.0.1:9000"
echo "  🐘 PostgreSQL:     localhost:5432"
echo "  📊 Redis:          localhost:6379"
echo "  🆔 Hydra Public:   http://127.0.0.1:4444"
echo "  🆔 Hydra Admin:    http://127.0.0.1:4445"

echo -e "\n${BLUE}📋 Quick Commands:${NC}"
echo "======================================"
echo "  • View logs: ./scripts/logs.sh [service-name]"
echo "  • Test deployment: ./scripts/test-deployment.sh"
echo "  • API docs: http://127.0.0.1:8080/v1/docs"
echo "  • GraphQL playground: http://127.0.0.1:8080/playground"

# Check for any failed pods
FAILED_PODS=$(kubectl get pods --all-namespaces -l app.kubernetes.io/part-of=naytife-platform --field-selector=status.phase!=Running --no-headers 2>/dev/null | wc -l)
if [ "$FAILED_PODS" -gt 0 ]; then
    echo -e "\n${YELLOW}⚠️  Warning: $FAILED_PODS pod(s) not in Running state${NC}"
    kubectl get pods --all-namespaces -l app.kubernetes.io/part-of=naytife-platform --field-selector=status.phase!=Running
fi
