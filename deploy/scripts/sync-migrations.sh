#!/bin/bash

# Migration files are now packaged in a dedicated Docker image
# This script is kept for backward compatibility and documentation purposes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}ℹ️  Migration files are now packaged in the backend-migrations Docker image${NC}"
echo -e "${BLUE}ℹ️  Source of truth: backend/internal/db/migrations/${NC}"
echo -e "${BLUE}ℹ️  Built automatically via: backend/migrations.Dockerfile${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_MIGRATIONS="$PROJECT_ROOT/../backend/internal/db/migrations"

echo -e "${YELLOW}� Current migration files in source:${NC}"

# Check if source directory exists
if [ ! -d "$BACKEND_MIGRATIONS" ]; then
    echo -e "${RED}❌ Source migrations directory not found: $BACKEND_MIGRATIONS${NC}"
    exit 1
fi

# Show current migration files
ls -la "$BACKEND_MIGRATIONS/"

echo ""
echo -e "${GREEN}✅ Migration files will be automatically included in the next image build${NC}"
echo -e "${BLUE}💡 To build the migration image locally: docker build -f backend/migrations.Dockerfile backend/${NC}"
echo -e "${BLUE}💡 To deploy: Use skaffold or your CI/CD pipeline${NC}"
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