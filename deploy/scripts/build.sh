#!/bin/bash

# Simple wrapper script for building all service images
# This provides an easy interface for the enhanced build system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_SCRIPT="$SCRIPT_DIR/build-images.sh"

echo -e "${BLUE}🏗️  Naytife Platform Image Builder${NC}"
echo ""

# Check if build script exists
if [ ! -f "$BUILD_SCRIPT" ]; then
    echo -e "${RED}❌ Build script not found: $BUILD_SCRIPT${NC}"
    exit 1
fi

# Show help if no arguments or help requested
if [ $# -eq 0 ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo -e "${YELLOW}Quick build commands:${NC}"
    echo -e "  ${GREEN}$0${NC}                    # Build only changed services"
    echo -e "  ${GREEN}$0 --force${NC}            # Force rebuild all services"
    echo -e "  ${GREEN}$0 --check${NC}            # Check which services need rebuilding"
    echo -e "  ${GREEN}$0 --service=backend${NC}  # Build only backend service"
    echo ""
    echo -e "${YELLOW}For advanced options, use: $BUILD_SCRIPT --help${NC}"
    echo ""
    
    if [ $# -eq 0 ]; then
        # Default behavior: check for changes and build if needed
        exec "$BUILD_SCRIPT"
    else
        exit 0
    fi
fi

# Convert simple flags to full arguments
args=()
for arg in "$@"; do
    case $arg in
        --check)
            args+=("--check-only")
            ;;
        *)
            args+=("$arg")
            ;;
    esac
done

# Execute the build script with processed arguments
exec "$BUILD_SCRIPT" "${args[@]}"
