#!/bin/bash

# Clean up obsolete template system services and files
# This script removes all traces of the old 4-service architecture

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧹 Cleaning up obsolete template system files${NC}"
echo "=============================================="

PROJECT_ROOT="/Users/erimebe/Development/commerce"

# List of obsolete services to remove
OBSOLETE_SERVICES=(
    "asset-manager"
    "template-builder" 
    "data-updater"
)

# Remove obsolete service directories
echo -e "\n${YELLOW}📁 Removing obsolete service directories...${NC}"
for service in "${OBSOLETE_SERVICES[@]}"; do
    SERVICE_DIR="$PROJECT_ROOT/services/$service"
    if [ -d "$SERVICE_DIR" ]; then
        echo -e "  ${RED}🗑️  Removing $SERVICE_DIR${NC}"
        rm -rf "$SERVICE_DIR"
    else
        echo -e "  ${GREEN}✅ $SERVICE_DIR already removed${NC}"
    fi
done

# Remove any remaining build artifacts
echo -e "\n${YELLOW}🔨 Removing build artifacts...${NC}"
for service in "${OBSOLETE_SERVICES[@]}"; do
    if [ -f "$PROJECT_ROOT/services/$service/$service" ]; then
        rm -f "$PROJECT_ROOT/services/$service/$service"
        echo -e "  ${GREEN}✅ Removed binary: $service${NC}"
    fi
done

# Check for any remaining references in scripts
echo -e "\n${YELLOW}🔍 Checking for remaining references...${NC}"
FOUND_REFS=0

for service in "${OBSOLETE_SERVICES[@]}"; do
    # Check k3s scripts
    if grep -r "$service" "$PROJECT_ROOT/k3s/scripts/" >/dev/null 2>&1; then
        echo -e "  ${RED}⚠️  Found references to $service in k3s scripts${NC}"
        FOUND_REFS=1
    fi
    
    # Check main scripts
    if grep -r "$service" "$PROJECT_ROOT/scripts/" >/dev/null 2>&1; then
        echo -e "  ${RED}⚠️  Found references to $service in main scripts${NC}"
        FOUND_REFS=1
    fi
done

if [ $FOUND_REFS -eq 0 ]; then
    echo -e "  ${GREEN}✅ No remaining references found${NC}"
fi

# List current template system structure
echo -e "\n${BLUE}📋 Current template system structure:${NC}"
echo "====================================="
echo "Services:"
ls -1 "$PROJECT_ROOT/services/" | grep -E "(template-registry|store-deployer)" | sed 's/^/  ✅ /'

echo ""
echo "Manifests:"
ls -1 "$PROJECT_ROOT/k3s/manifests/08-template-system/" | sed 's/^/  ✅ /'

echo ""
echo "Scripts:"
ls -1 "$PROJECT_ROOT/scripts/" | grep -E "(upload-template|test-template)" | sed 's/^/  ✅ /'

echo -e "\n${GREEN}🎉 Cleanup completed successfully!${NC}"
echo -e "${GREEN}✨ Template system is now using the new 2-service architecture${NC}"
