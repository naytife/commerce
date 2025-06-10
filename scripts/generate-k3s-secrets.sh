#!/bin/bash

# Generate k3s secrets from .env file
# This script reads from .env and creates/updates k3s secret files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔑 Generating k3s secrets from .env file${NC}"

# Function to load environment variables from .env file
load_env() {
    if [[ -f .env ]]; then
        echo -e "${GREEN}Loading environment variables from .env file...${NC}"
        export $(grep -v '^#' .env | xargs)
    else
        echo -e "${RED}Error: .env file not found${NC}"
        echo "Please copy .env.example to .env and configure your values."
        exit 1
    fi
}

# Load environment variables
load_env

# Generate backend secrets
echo -e "${YELLOW}Generating backend secrets...${NC}"
cat > k3s/backend/secrets.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: backend-secrets
  namespace: backend
type: Opaque
stringData:
  # Database configuration
  DATABASE_URL: "${DATABASE_URL/localhost/postgresql.database.svc.cluster.local}"
  
  # Redis configuration  
  REDIS_URL: "${REDIS_URL/localhost/redis-master.database.svc.cluster.local}"
  
  # Stripe configuration
  STRIPE_WEBHOOK_SECRET: "${STRIPE_WEBHOOK_SECRET}"
  STRIPE_SECRET_KEY: "${STRIPE_SECRET_KEY}"
  STRIPE_PUBLISHABLE_KEY: "${STRIPE_PUBLISHABLE_KEY}"
  
  # JWT configuration
  JWT_SECRET: "${JWT_SECRET}"
  
  # Server configuration
  ENVIRONMENT: "${ENVIRONMENT:-development}"
  ALLOWED_ORIGINS: "${ALLOWED_ORIGINS:-*}"
  
  # OAuth2 configuration
  OAUTH2_CLIENT_ID: "${OAUTH2_CLIENT_ID}"
  OAUTH2_CLIENT_SECRET: "${OAUTH2_CLIENT_SECRET}"
  OAUTH2_REDIRECT_URL: "${OAUTH2_REDIRECT_URL}"
EOF

# Generate auth secrets
echo -e "${YELLOW}Generating auth secrets...${NC}"
cat > k3s/auth/secrets.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: auth-secrets
  namespace: auth
type: Opaque
stringData:
  # Database configuration for Hydra
  DATABASE_URL: "${DATABASE_URL/localhost/postgresql.database.svc.cluster.local}"
  
  # OAuth2 configuration
  OAUTH2_CLIENT_ID: "${OAUTH2_CLIENT_ID}"
  OAUTH2_CLIENT_SECRET: "${OAUTH2_CLIENT_SECRET}"
  OAUTH2_REDIRECT_URL: "${OAUTH2_REDIRECT_URL}"
  
  # JWT configuration
  JWT_SECRET: "${JWT_SECRET}"
  
  # Hydra system secret
  HYDRA_SYSTEM_SECRET: "${HYDRA_SYSTEM_SECRET:-dev-hydra-system-secret}"
EOF

# Generate cloud-build secrets
echo -e "${YELLOW}Generating cloud-build secrets...${NC}"
cat > k3s/cloud-build/secrets.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: cloud-build-secrets
  namespace: cloud-build
type: Opaque
stringData:
  # Redis configuration
  REDIS_URL: "${REDIS_URL/localhost/redis-master.database.svc.cluster.local}"
  
  # Cloudflare R2 configuration
  R2_ACCESS_KEY: "${R2_ACCESS_KEY}"
  R2_SECRET_KEY: "${R2_SECRET_KEY}"
  R2_ENDPOINT: "${R2_ENDPOINT}"
  R2_BUCKET_REGION: "${R2_BUCKET_REGION:-auto}"
  CLOUDFLARE_ACCOUNT_ID: "${CLOUDFLARE_ACCOUNT_ID}"
  CLOUDFLARE_ACCESS_KEY_ID: "${CLOUDFLARE_ACCESS_KEY_ID:-${R2_ACCESS_KEY}}"
  CLOUDFLARE_SECRET_ACCESS_KEY: "${CLOUDFLARE_SECRET_ACCESS_KEY:-${R2_SECRET_KEY}}"
  
  # Build configuration
  TEMPLATES_DIR: "${TEMPLATES_DIR:-/app/templates}"
  BUILD_DIR: "${BUILD_DIR:-/app/built_sites}"
EOF

echo -e "${GREEN}✅ All k3s secrets generated successfully!${NC}"
echo ""
echo -e "${BLUE}Generated files:${NC}"
echo "  - k3s/backend/secrets.yaml"
echo "  - k3s/auth/secrets.yaml" 
echo "  - k3s/cloud-build/secrets.yaml"
echo ""
echo -e "${YELLOW}Note: Remember to apply these secrets after deploying namespaces:${NC}"
echo "  kubectl apply -f k3s/backend/secrets.yaml"
echo "  kubectl apply -f k3s/auth/secrets.yaml"
echo "  kubectl apply -f k3s/cloud-build/secrets.yaml"
