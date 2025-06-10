#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SERVICE_NAME=${1}
NAMESPACE=${2}
FOLLOW=${3:-"false"}

# Service mappings (service-name -> namespace)
declare -A SERVICE_NAMESPACES=(
    ["postgres"]="naytife"
    ["redis"]="naytife"
    ["backend"]="naytife"
    ["hydra"]="naytife-auth"
    ["hydra-migrate"]="naytife-auth"
    ["oathkeeper"]="naytife-auth"
    ["auth-handler"]="naytife-auth"
    ["cloud-build"]="naytife-build"
)

show_usage() {
    echo -e "${BLUE}📋 Naytife Log Viewer${NC}"
    echo "===================="
    echo "Usage: $0 [service-name] [namespace] [follow]"
    echo ""
    echo -e "${YELLOW}Available services:${NC}"
    for service in "${!SERVICE_NAMESPACES[@]}"; do
        echo "  • $service (${SERVICE_NAMESPACES[$service]})"
    done
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 backend                    # View backend logs"
    echo "  $0 backend naytife true      # Follow backend logs"
    echo "  $0 hydra                     # View Hydra logs"
    echo "  $0                           # Show all pod logs"
}

if [ -z "$SERVICE_NAME" ]; then
    show_usage
    echo ""
    echo -e "${BLUE}📊 All Pod Logs (last 50 lines each):${NC}"
    echo "====================================="
    
    # Show logs from all naytife pods
    for ns in naytife naytife-auth naytife-build; do
        echo -e "\n${YELLOW}Namespace: $ns${NC}"
        echo "------------------------"
        
        pods=$(kubectl get pods -n "$ns" -l app.kubernetes.io/part-of=naytife-platform --no-headers -o custom-columns=":metadata.name" 2>/dev/null)
        
        if [ -z "$pods" ]; then
            echo -e "${RED}No pods found in namespace $ns${NC}"
            continue
        fi
        
        for pod in $pods; do
            echo -e "\n${BLUE}📦 Pod: $pod${NC}"
            echo "-------------------"
            kubectl logs -n "$ns" "$pod" --tail=20 2>/dev/null || echo -e "${RED}Failed to get logs for $pod${NC}"
        done
    done
    exit 0
fi

# Determine namespace if not provided
if [ -z "$NAMESPACE" ]; then
    if [ -n "${SERVICE_NAMESPACES[$SERVICE_NAME]}" ]; then
        NAMESPACE="${SERVICE_NAMESPACES[$SERVICE_NAME]}"
    else
        echo -e "${RED}❌ Unknown service: $SERVICE_NAME${NC}"
        echo -e "${YELLOW}Available services: ${!SERVICE_NAMESPACES[*]}${NC}"
        exit 1
    fi
fi

echo -e "${BLUE}📋 Viewing logs for $SERVICE_NAME in namespace $NAMESPACE${NC}"
echo "=============================================="

# Check if deployment exists
if ! kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE" >/dev/null 2>&1; then
    echo -e "${RED}❌ Deployment $SERVICE_NAME not found in namespace $NAMESPACE${NC}"
    
    # Try to find pods with similar names
    echo -e "\n${YELLOW}Looking for pods with similar names...${NC}"
    kubectl get pods -n "$NAMESPACE" | grep -i "$SERVICE_NAME" || echo "No similar pods found"
    exit 1
fi

# Get pod name
POD_NAME=$(kubectl get pods -n "$NAMESPACE" -l app="$SERVICE_NAME" --no-headers -o custom-columns=":metadata.name" | head -1)

if [ -z "$POD_NAME" ]; then
    echo -e "${RED}❌ No running pods found for $SERVICE_NAME${NC}"
    echo -e "\n${YELLOW}Deployment status:${NC}"
    kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE"
    exit 1
fi

echo -e "${GREEN}✅ Found pod: $POD_NAME${NC}"

# Show logs
if [ "$FOLLOW" = "true" ] || [ "$3" = "-f" ] || [ "$3" = "--follow" ]; then
    echo -e "${BLUE}📡 Following logs for $POD_NAME...${NC}"
    echo "Press Ctrl+C to stop"
    echo "========================"
    kubectl logs -n "$NAMESPACE" "$POD_NAME" -f
else
    echo -e "${BLUE}📄 Recent logs for $POD_NAME:${NC}"
    echo "========================"
    kubectl logs -n "$NAMESPACE" "$POD_NAME" --tail=100
fi

echo -e "\n${BLUE}💡 Tip:${NC} Use '$0 $SERVICE_NAME $NAMESPACE true' to follow logs in real-time"
