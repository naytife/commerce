#!/bin/bash

# Script to create Hydra OAuth2 clients for development
# This script creates the necessary OAuth2 clients for:
# - Dashboard application (SvelteKit)
# - Swagger UI documentation
# - Oathkeeper proxy authentication

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔑 Creating Hydra OAuth2 Clients${NC}"
echo "========================================"

# Wait for Hydra to be ready
echo -e "${YELLOW}⏳ Waiting for Hydra admin to be ready...${NC}"
while ! kubectl exec -n naytife-auth deploy/local-hydra -- hydra list oauth2-clients --endpoint http://localhost:4445 >/dev/null 2>&1; do
    echo "   Waiting for Hydra admin..."
    sleep 2
done
echo -e "${GREEN}✅ Hydra admin is ready${NC}"

# Function to create a client using Hydra admin API
create_client_with_id() {
    local client_id="$1"
    local client_secret="$2"
    local client_name="$3"
    local redirect_uris="$4"
    local grant_types="$5"
    local response_types="$6"
    local scopes="$7"
    local token_endpoint_auth_method="$8"

    echo -e "${YELLOW}📝 Creating client: ${client_name} (ID: ${client_id})${NC}"
    
    # Check if client already exists and delete it
    if kubectl exec -n naytife-auth deploy/local-hydra -- hydra get oauth2-client --endpoint http://localhost:4445 "$client_id" >/dev/null 2>&1; then
        echo -e "${YELLOW}   Client already exists, deleting...${NC}"
        kubectl exec -n naytife-auth deploy/local-hydra -- hydra delete oauth2-client --endpoint http://localhost:4445 "$client_id" >/dev/null 2>&1
    fi
    
    # Create JSON payload
    local redirect_uris_json=""
    if [ -n "$redirect_uris" ]; then
        redirect_uris_json="\"$redirect_uris\""
    fi
    
    # Convert comma-separated values to JSON arrays
    local grant_types_json=$(echo "[$grant_types]" | sed 's/,/","/g' | sed 's/\[/["/' | sed 's/\]/"]/')
    local response_types_json=""
    if [ -n "$response_types" ]; then
        response_types_json=$(echo "[$response_types]" | sed 's/,/","/g' | sed 's/\[/["/' | sed 's/\]/"]/')
    fi
    
    # Scopes should be space-separated for Hydra API (pass as-is)
    local scopes_formatted="$scopes"
    
    # Create the client using a temporary pod with curl
    local json_payload
    if [ -n "$redirect_uris" ]; then
        json_payload="{
            \"client_id\": \"$client_id\",
            \"client_secret\": \"$client_secret\",
            \"client_name\": \"$client_name\",
            \"grant_types\": $grant_types_json,
            \"response_types\": $response_types_json,
            \"scope\": \"$scopes_formatted\",
            \"redirect_uris\": [$redirect_uris_json],
            \"token_endpoint_auth_method\": \"$token_endpoint_auth_method\"
        }"
    else
        json_payload="{
            \"client_id\": \"$client_id\",
            \"client_secret\": \"$client_secret\",
            \"client_name\": \"$client_name\",
            \"grant_types\": $grant_types_json,
            \"scope\": \"$scopes_formatted\",
            \"token_endpoint_auth_method\": \"$token_endpoint_auth_method\"
        }"
    fi
    
    kubectl run temp-curl --rm -i --image=curlimages/curl:latest --restart=Never -- \
        curl -s -X POST http://local-hydra-admin.naytife-auth:4445/admin/clients \
        -H "Content-Type: application/json" \
        -d "$json_payload" > /tmp/create_result.json 2>/dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}   ✅ Created: ${client_name}${NC}"
    else
        echo -e "${RED}   ❌ Failed to create: ${client_name}${NC}"
        return 1
    fi
}

# 1. Dashboard Client (SvelteKit)
echo ""
create_client_with_id \
    "4b41cd38-43ed-4e3a-9a88-bd384af21732" \
    "fbOoeUd9fEiw6LM~TWhg70zhTo" \
    "Dashboard Application" \
    "http://localhost:5173/auth/callback/hydra" \
    "authorization_code,refresh_token" \
    "code" \
    "openid offline_access hydra.openid introspect" \
    "client_secret_basic"

# 2. Swagger UI Client
echo ""
create_client_with_id \
    "d39beaaa-9c53-48e7-b82a-37ff52127473" \
    "-tzS7OuCyHjTZUxtfx5TxGR1f." \
    "Swagger UI Documentation" \
    "http://127.0.0.1:8080/v1/docs/oauth2-redirect.html" \
    "authorization_code,refresh_token" \
    "code" \
    "openid offline profile email offline_access" \
    "client_secret_basic"

# 3. Oathkeeper Client (for introspection)
echo ""
create_client_with_id \
    "761506cc-511f-411c-bf31-752efd8063b3" \
    "z.WE65SP5o0oXDYJwLdyoqYUuN" \
    "Oathkeeper Proxy" \
    "" \
    "client_credentials" \
    "" \
    "introspect" \
    "client_secret_basic"

echo ""
echo -e "${GREEN}🎉 All OAuth2 clients created successfully!${NC}"
echo ""
echo -e "${BLUE}📋 Client Summary:${NC}"
echo "=================================================="
echo -e "${YELLOW}Dashboard Application:${NC}"
echo "  Client ID:     4b41cd38-43ed-4e3a-9a88-bd384af21732"
echo "  Redirect URI:  http://localhost:5173/auth/callback/hydra"
echo ""
echo -e "${YELLOW}Swagger UI Documentation:${NC}"
echo "  Client ID:     d39beaaa-9c53-48e7-b82a-37ff52127473"
echo "  Redirect URI:  http://127.0.0.1:8080/v1/docs/oauth2-redirect.html"
echo ""
echo -e "${YELLOW}Oathkeeper Proxy:${NC}"
echo "  Client ID:     761506cc-511f-411c-bf31-752efd8063b3"
echo "  Purpose:       Token introspection"
echo ""
echo -e "${GREEN}✅ Ready for development!${NC}"
