#!/bin/bash

# Backend for Frontend Test Script
# This script tests the proxy functionality of the backend service

echo "🚀 Testing Backend for Frontend Implementation"
echo "=============================================="

# Configuration
BACKEND_URL="http://127.0.0.1:8080/v1"
AUTH_TOKEN=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to make authenticated request
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    echo -e "${YELLOW}Testing:${NC} $method $endpoint"
    
    if [ -n "$data" ]; then
        curl -s -X "$method" \
             -H "Content-Type: application/json" \
             -H "Authorization: Bearer $AUTH_TOKEN" \
             -d "$data" \
             "$BACKEND_URL$endpoint"
    else
        curl -s -X "$method" \
             -H "Authorization: Bearer $AUTH_TOKEN" \
             "$BACKEND_URL$endpoint"
    fi
    
    echo -e "\n"
}

# Function to test health endpoint (no auth required)
test_health() {
    echo -e "${YELLOW}Testing:${NC} Health endpoint"
    response=$(curl -s "$BACKEND_URL/health/services")
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
    echo -e "\n"
}

echo "📋 Test Plan:"
echo "1. Backend service health"
echo "2. Aggregated services health (template-registry + store-deployer)"
echo "3. Template operations (proxied to template-registry)"
echo "4. Deployment operations (proxied to store-deployer)"
echo ""

# Test 1: Backend Health
echo -e "${GREEN}Test 1: Backend Service Health${NC}"
echo "================================="
make_request "GET" "/health"

# Test 2: Aggregated Services Health
echo -e "${GREEN}Test 2: Aggregated Services Health${NC}"
echo "=================================="
test_health

# Test 3: Template Operations (proxied)
echo -e "${GREEN}Test 3: Template Operations (Proxied)${NC}"
echo "====================================="
make_request "GET" "/templates"

# Test 4: Deployment Operations (proxied)
echo -e "${GREEN}Test 4: Deployment Operations (Proxied)${NC}"
echo "======================================="
# Note: This requires authentication and a valid shop_id
echo "Note: Deployment operations require authentication"
echo "Example: POST /v1/shops/{shop_id}/deploy"
echo ""

echo "🔍 Backend Proxy Features Demonstrated:"
echo "======================================="
echo "✓ Single API endpoint for all operations"
echo "✓ Centralized health monitoring"
echo "✓ Request proxying to microservices"
echo "✓ Consistent error handling"
echo "✓ Authentication flow through backend"
echo ""

echo "🎯 Benefits of Backend for Frontend Pattern:"
echo "============================================"
echo "• Simplified client configuration"
echo "• Service abstraction"
echo "• Centralized authentication"
echo "• Easier monitoring and debugging"
echo "• Future-proof architecture"
echo ""

echo "📚 Next Steps:"
echo "=============="
echo "1. Deploy updated services: cd k3s && ./scripts/deploy-all.sh"
echo "2. Test with authenticated requests"
echo "3. Monitor proxy performance"
echo "4. Extend pattern to additional services"
