#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Skaffold Development Utility${NC}"
echo -e "${YELLOW}Available commands:${NC}"

show_help() {
    echo -e "${GREEN}Usage: $0 <command> [options]${NC}"
    echo ""
    echo -e "${YELLOW}Commands:${NC}"
    echo -e "  ${GREEN}build${NC}        - Build all images without deploying"
    echo -e "  ${GREEN}build-smart${NC}  - Build only changed images (with change detection)"
    echo -e "  ${GREEN}render${NC}       - Render Kubernetes manifests without deploying"
    echo -e "  ${GREEN}deploy${NC}       - Deploy using pre-built images"
    echo -e "  ${GREEN}delete${NC}       - Delete all deployed resources"
    echo -e "  ${GREEN}status${NC}       - Show status of deployed resources"
    echo -e "  ${GREEN}logs${NC}         - Tail logs from all services"
    echo -e "  ${GREEN}debug${NC}        - Start in debug mode"
    echo -e "  ${GREEN}validate${NC}     - Validate Skaffold configuration"
    echo -e "  ${GREEN}clean${NC}        - Clean up Docker images and build cache"
    echo -e "  ${GREEN}check${NC}        - Check which services need rebuilding"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo -e "  $0 build"
    echo -e "  $0 build-smart"
    echo -e "  $0 build-smart --service=backend"
    echo -e "  $0 render"
    echo -e "  $0 deploy"
    echo -e "  $0 logs backend"
    echo -e "  $0 debug"
    echo -e "  $0 check"
}

# Get script directory for locating build script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check if skaffold is available
if ! command -v skaffold &> /dev/null; then
    echo -e "${RED}❌ Skaffold not installed. Please install Skaffold first.${NC}"
    exit 1
fi

COMMAND=${1:-help}
shift || true  # Remove first argument, keep the rest for passing to build script

case $COMMAND in
    "build")
        echo -e "${BLUE}🔨 Building all images with Skaffold...${NC}"
        skaffold build --profile=local
        ;;
    
    "build-smart")
        echo -e "${BLUE}🧠 Building images with change detection...${NC}"
        if [ -f "$SCRIPT_DIR/build-images.sh" ]; then
            "$SCRIPT_DIR/build-images.sh" "$@"
        else
            echo -e "${RED}❌ build-images.sh not found. Please ensure it exists.${NC}"
            exit 1
        fi
        ;;
    
    "check")
        echo -e "${BLUE}🔍 Checking which services need rebuilding...${NC}"
        if [ -f "$SCRIPT_DIR/build-images.sh" ]; then
            "$SCRIPT_DIR/build-images.sh" --check-only "$@"
        else
            echo -e "${RED}❌ build-images.sh not found. Please ensure it exists.${NC}"
            exit 1
        fi
        ;;
    
    "render")
        echo -e "${BLUE}📄 Rendering Kubernetes manifests...${NC}"
        skaffold render --profile=local
        ;;
    
    "deploy")
        echo -e "${BLUE}🚀 Deploying using pre-built images...${NC}"
        skaffold deploy --profile=local
        ;;
    
    "delete")
        echo -e "${YELLOW}🗑️  Deleting all deployed resources...${NC}"
        read -p "Are you sure you want to delete all resources? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            skaffold delete --profile=local
            echo -e "${GREEN}✅ Resources deleted${NC}"
        else
            echo -e "${YELLOW}❌ Deletion cancelled${NC}"
        fi
        ;;
    
    "status")
        echo -e "${BLUE}📊 Checking deployment status...${NC}"
        echo -e "${YELLOW}Namespace: naytife-local${NC}"
        kubectl get pods -n naytife-local
        echo ""
        echo -e "${YELLOW}Services:${NC}"
        kubectl get services -n naytife-local
        ;;
    
    "logs")
        SERVICE=${2:-""}
        if [ -n "$SERVICE" ]; then
            echo -e "${BLUE}📋 Tailing logs for service: $SERVICE${NC}"
            kubectl logs -f -n naytife-local deployment/$SERVICE
        else
            echo -e "${BLUE}📋 Tailing logs for all services...${NC}"
            echo -e "${YELLOW}Available services: backend, auth-handler, store-deployer, template-registry${NC}"
            echo -e "${YELLOW}Use: $0 logs <service-name>${NC}"
        fi
        ;;
    
    "debug")
        echo -e "${BLUE}🐛 Starting debug mode...${NC}"
        skaffold dev --profile=debug --port-forward --cleanup
        ;;
    
    "validate")
        echo -e "${BLUE}✅ Validating Skaffold configuration...${NC}"
        skaffold diagnose
        echo ""
        echo -e "${BLUE}🔍 Checking configuration syntax...${NC}"
        skaffold config list
        ;;
    
    "clean")
        echo -e "${YELLOW}🧹 Cleaning up Docker images and build cache...${NC}"
        read -p "This will remove unused Docker images and build cache. Continue? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}🧹 Cleaning Docker system...${NC}"
            docker system prune -f
            echo -e "${BLUE}🧹 Cleaning build cache...${NC}"
            docker builder prune -f
            echo -e "${GREEN}✅ Cleanup completed${NC}"
        else
            echo -e "${YELLOW}❌ Cleanup cancelled${NC}"
        fi
        ;;
    
    "help"|*)
        show_help
        ;;
esac
