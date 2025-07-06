#!/bin/bash

# Sync migration files from backend to deploy directory
# This ensures the migration files in deploy/ are always up-to-date

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_MIGRATIONS="$PROJECT_ROOT/../backend/internal/db/migrations"
DEPLOY_MIGRATIONS="$SCRIPT_DIR/../base/backend/migrations"

echo -e "${YELLOW}🔄 Syncing migration files...${NC}"

# Check if source directory exists
if [ ! -d "$BACKEND_MIGRATIONS" ]; then
    echo -e "${RED}❌ Source migrations directory not found: $BACKEND_MIGRATIONS${NC}"
    exit 1
fi

# Create deploy migrations directory if it doesn't exist
mkdir -p "$DEPLOY_MIGRATIONS"

# Copy migration files
echo "📁 Copying migration files from backend to deploy..."
cp "$BACKEND_MIGRATIONS"/*.sql "$DEPLOY_MIGRATIONS/" 2>/dev/null || true
cp "$BACKEND_MIGRATIONS"/atlas.sum "$DEPLOY_MIGRATIONS/" 2>/dev/null || true

# Check if we have migration files
if [ ! -f "$DEPLOY_MIGRATIONS/atlas.sum" ]; then
    echo -e "${RED}❌ No migration files found in $BACKEND_MIGRATIONS${NC}"
    exit 1
fi

echo "📋 Migration files synced:"
ls -la "$DEPLOY_MIGRATIONS/"

# Validate the kustomization
echo "🔍 Validating Kustomize configuration..."
if kustomize build "$SCRIPT_DIR/../base/backend" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Kustomize validation passed${NC}"
else
    echo -e "${RED}❌ Kustomize validation failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Migration sync complete!${NC}"
echo ""
echo "Migration files are now in sync with your backend directory."
echo "The ConfigMapGenerator will automatically include these files in your deployment."
echo ""
echo "Next steps:"
echo "  1. Deploy with: kubectl apply -k deploy/overlays/local"
echo "  2. Check migration job: kubectl get jobs -n naytife"
echo "  3. View migration logs: kubectl logs job/backend-migrate -n naytife" 