#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Naytife Platform with Skaffold...${NC}"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo -e "${YELLOW}🔍 Checking prerequisites...${NC}"

if ! command_exists kubectl; then
    echo -e "${RED}❌ kubectl not installed. Please install kubectl first.${NC}"
    exit 1
fi

if ! command_exists skaffold; then
    echo -e "${RED}❌ Skaffold not installed. Please install Skaffold first.${NC}"
    echo -e "${YELLOW}   Install with: brew install skaffold${NC}"
    exit 1
fi

# Check if k3s is running
echo -e "${YELLOW}🔍 Checking k3s cluster connectivity...${NC}"
if ! kubectl get nodes &> /dev/null; then
    echo -e "${YELLOW}⚠️  k3s cluster not accessible. Starting k3s...${NC}"
    if [ -f "./k3s/scripts/create-cluster.sh" ]; then
        ./k3s/scripts/create-cluster.sh
    else
        echo -e "${RED}❌ k3s setup script not found. Please ensure k3s is running.${NC}"
        exit 1
    fi
fi

# Verify k3s is accessible
if ! kubectl get nodes &> /dev/null; then
    echo -e "${RED}❌ Cannot connect to k3s cluster. Please check your setup.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ k3s cluster is accessible${NC}"

# Check if namespace exists, create if not
NAMESPACE="naytife-local"
if ! kubectl get namespace $NAMESPACE &> /dev/null; then
    echo -e "${YELLOW}🔧 Creating namespace: $NAMESPACE${NC}"
    kubectl create namespace $NAMESPACE
else
    echo -e "${GREEN}✅ Namespace $NAMESPACE exists${NC}"
fi

# Determine deployment mode
MODE=${1:-dev}
PROFILE="local"
CLEANUP=${2:-true}

case $MODE in
    "dev")
        echo -e "${BLUE}🔧 Starting development mode with file watching...${NC}"
        SKAFFOLD_CMD="dev"
        ;;
    "run")
        echo -e "${BLUE}🔧 Building and deploying (one-time)...${NC}"
        SKAFFOLD_CMD="run"
        ;;
    "debug")
        echo -e "${BLUE}🔧 Starting debug mode...${NC}"
        SKAFFOLD_CMD="dev"
        PROFILE="debug"
        ;;
    "fast")
        echo -e "${BLUE}🔧 Starting fast build mode...${NC}"
        SKAFFOLD_CMD="dev"
        PROFILE="fast"
        ;;
    *)
        echo -e "${RED}❌ Invalid mode: $MODE${NC}"
        echo -e "${YELLOW}   Valid modes: dev, run, debug, fast${NC}"
        exit 1
        ;;
esac

# Build command
SKAFFOLD_ARGS="--profile=$PROFILE"
if [ "$CLEANUP" = "true" ] && [ "$SKAFFOLD_CMD" = "dev" ]; then
    SKAFFOLD_ARGS="$SKAFFOLD_ARGS --cleanup"
fi

if [ "$SKAFFOLD_CMD" = "dev" ]; then
    SKAFFOLD_ARGS="$SKAFFOLD_ARGS --port-forward"
fi

echo -e "${BLUE}🔧 Deploying with Skaffold...${NC}"
echo -e "${YELLOW}   Command: skaffold $SKAFFOLD_CMD $SKAFFOLD_ARGS${NC}"

# Run Skaffold
if skaffold $SKAFFOLD_CMD $SKAFFOLD_ARGS; then
    echo -e "${GREEN}✅ Naytife Platform is running!${NC}"
    echo -e "${BLUE}📡 Access points:${NC}"
    echo -e "${GREEN}   Backend: http://localhost:8000${NC}"
    echo -e "${GREEN}   Auth Handler: http://localhost:3000${NC}"
    echo -e "${GREEN}   Oathkeeper: http://localhost:4456${NC}"
    echo -e "${GREEN}   Hydra Admin: http://localhost:4445${NC}"
    echo -e "${GREEN}   PostgreSQL: localhost:5432${NC}"
    echo -e "${YELLOW}   Use Ctrl+C to stop the development environment${NC}"
else
    echo -e "${RED}❌ Skaffold deployment failed${NC}"
    exit 1
fi
