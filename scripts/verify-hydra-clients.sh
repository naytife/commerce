#!/bin/bash

# Script to verify Hydra OAuth2 clients are properly configured

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Verifying Hydra OAuth2 Clients${NC}"
echo "====================================="

# Check if Hydra is accessible
if ! kubectl exec -n naytife-auth deploy/hydra -- hydra list oauth2-clients --endpoint http://localhost:4445 >/dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to Hydra admin${NC}"
    echo "   Make sure the k3s cluster is running and Hydra is deployed"
    exit 1
fi

echo -e "${GREEN}✅ Hydra admin is accessible${NC}"
echo ""

# Function to check if client exists
check_client() {
    local client_id="$1"
    local client_name="$2"
    
    echo -n "🔑 Checking $client_name: "
    
    if kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 "$client_id" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Found${NC}"
        return 0
    else
        echo -e "${RED}❌ Missing${NC}"
        return 1
    fi
}

# Check all required clients
missing_clients=0

check_client "4b41cd38-43ed-4e3a-9a88-bd384af21732" "Dashboard Application" || ((missing_clients++))
check_client "d39beaaa-9c53-48e7-b82a-37ff52127473" "Swagger UI Documentation" || ((missing_clients++))
check_client "761506cc-511f-411c-bf31-752efd8063b3" "Oathkeeper Proxy" || ((missing_clients++))

echo ""

if [ $missing_clients -eq 0 ]; then
    echo -e "${GREEN}🎉 All OAuth2 clients are properly configured!${NC}"
    echo ""
    echo -e "${BLUE}📋 Client Details:${NC}"
    echo "=================================================="
    
    # Show detailed client information
    echo -e "${YELLOW}Dashboard Application:${NC}"
    kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 "4b41cd38-43ed-4e3a-9a88-bd384af21732" 2>/dev/null | grep -E "(Client ID|Name|Redirect URIs|Grant Types|Scopes)" || echo "   Details available via: kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 4b41cd38-43ed-4e3a-9a88-bd384af21732"
    
    echo ""
    echo -e "${YELLOW}Swagger UI Documentation:${NC}"
    kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 "d39beaaa-9c53-48e7-b82a-37ff52127473" 2>/dev/null | grep -E "(Client ID|Name|Redirect URIs|Grant Types|Scopes)" || echo "   Details available via: kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 d39beaaa-9c53-48e7-b82a-37ff52127473"
    
    echo ""
    echo -e "${YELLOW}Oathkeeper Proxy:${NC}"
    kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 "761506cc-511f-411c-bf31-752efd8063b3" 2>/dev/null | grep -E "(Client ID|Name|Grant Types|Scopes)" || echo "   Details available via: kubectl exec -n naytife-auth deploy/hydra -- hydra get oauth2-client --endpoint http://localhost:4445 761506cc-511f-411c-bf31-752efd8063b3"
    
    exit 0
else
    echo -e "${RED}❌ $missing_clients client(s) missing${NC}"
    echo ""
    echo -e "${YELLOW}💡 To create missing clients, run:${NC}"
    echo "   ./scripts/create-hydra-clients.sh"
    exit 1
fi
