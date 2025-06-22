#!/bin/bash

# Script to display OAuth2 client credentials for development
# Use this for quick reference when configuring applications

# Colors for output
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔑 Hydra OAuth2 Client Credentials${NC}"
echo "===================================="
echo ""

echo -e "${YELLOW}Dashboard Application (SvelteKit):${NC}"
echo "  Client ID:     4b41cd38-43ed-4e3a-9a88-bd384af21732"
echo "  Client Secret: fbOoeUd9fEiw6LM~TWhg70zhTo"
echo "  Redirect URI:  http://localhost:5173/auth/callback/hydra"
echo "  Scopes:        openid offline hydra.openid introspect"
echo ""

echo -e "${YELLOW}Swagger UI Documentation:${NC}"
echo "  Client ID:     d39beaaa-9c53-48e7-b82a-37ff52127473"
echo "  Client Secret: -tzS7OuCyHjTZUxtfx5TxGR1f."
echo "  Redirect URI:  http://127.0.0.1:8080/v1/docs/oauth2-redirect.html"
echo "  Scopes:        openid offline profile email offline_access"
echo ""

echo -e "${YELLOW}Oathkeeper Proxy (Introspection):${NC}"
echo "  Client ID:     761506cc-511f-411c-bf31-752efd8063b3"
echo "  Client Secret: z.WE65SP5o0oXDYJwLdyoqYUuN"
echo "  Grant Type:    client_credentials"
echo "  Scopes:        introspect"
echo ""

echo -e "${GREEN}💡 Usage:${NC}"
echo "• Run './scripts/create-hydra-clients.sh' to create these clients"
echo "• Run './scripts/verify-hydra-clients.sh' to check if they exist"
echo "• Run './scripts/setup-dev-environment.sh' for complete setup"
