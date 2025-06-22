#!/bin/bash

# Test Bucket Separation Script
# This script tests that the bucket separation is working correctly

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing Bucket Separation${NC}"
echo -e "${BLUE}=============================${NC}"

# Test Kubernetes services
test_kubernetes_services() {
    echo -e "\n${YELLOW}☸️  Testing Kubernetes Services${NC}"
    
    # Check if services are running
    echo -e "${BLUE}Checking service status...${NC}"
    kubectl get pods -n naytife -l app=template-registry
    kubectl get pods -n naytife -l app=store-deployer
    
    # Check service logs for bucket configuration
    echo -e "\n${BLUE}Checking template-registry bucket configuration...${NC}"
    template_logs=$(kubectl logs deployment/template-registry -n naytife --tail=20)
    if echo "$template_logs" | grep -q "templates"; then
        echo -e "${GREEN}✅ Template-registry appears to be using templates bucket${NC}"
    else
        echo -e "${RED}❌ Template-registry bucket configuration unclear${NC}"
        echo "Recent logs:"
        echo "$template_logs"
    fi
    
    echo -e "\n${BLUE}Checking store-deployer bucket configuration...${NC}"
    store_logs=$(kubectl logs deployment/store-deployer -n naytife --tail=20)
    if echo "$store_logs" | grep -q "stores"; then
        echo -e "${GREEN}✅ Store-deployer appears to be using stores bucket${NC}"
    else
        echo -e "${RED}❌ Store-deployer bucket configuration unclear${NC}"
        echo "Recent logs:"
        echo "$store_logs"
    fi
}

# Test API endpoints
test_api_endpoints() {
    echo -e "\n${YELLOW}🌐 Testing API Endpoints${NC}"
    
    # Test template registry health
    echo -e "${BLUE}Testing template-registry health...${NC}"
    if curl -s http://127.0.0.1:9001/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Template-registry is responding${NC}"
    else
        echo -e "${RED}❌ Template-registry is not responding${NC}"
    fi
    
    # Test store deployer health
    echo -e "${BLUE}Testing store-deployer health...${NC}"
    if curl -s http://127.0.0.1:9003/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Store-deployer is responding${NC}"
    else
        echo -e "${RED}❌ Store-deployer is not responding${NC}"
    fi
    
    # Test templates endpoint
    echo -e "${BLUE}Testing templates endpoint...${NC}"
    if curl -s http://127.0.0.1:9001/templates > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Templates endpoint is working${NC}"
    else
        echo -e "${RED}❌ Templates endpoint is not working${NC}"
    fi
}

# Test bucket access with rclone (if available)
test_bucket_access() {
    echo -e "\n${YELLOW}📦 Testing Bucket Access${NC}"
    
    if command -v rclone &> /dev/null; then
        echo -e "${BLUE}Testing templates bucket access...${NC}"
        if rclone ls r2:templates/ > /dev/null 2>&1; then
            template_count=$(rclone ls r2:templates/ | wc -l)
            echo -e "${GREEN}✅ Templates bucket accessible ($template_count files)${NC}"
        else
            echo -e "${RED}❌ Templates bucket not accessible${NC}"
        fi
        
        echo -e "${BLUE}Testing stores bucket access...${NC}"
        if rclone ls r2:stores/ > /dev/null 2>&1; then
            stores_count=$(rclone ls r2:stores/ | wc -l)
            echo -e "${GREEN}✅ Stores bucket accessible ($stores_count files)${NC}"
        else
            echo -e "${RED}❌ Stores bucket not accessible${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  rclone not available, skipping bucket access tests${NC}"
    fi
}

# Check environment variables
check_environment() {
    echo -e "\n${YELLOW}🔧 Checking Environment Configuration${NC}"
    
    # Check Kubernetes secrets
    echo -e "${BLUE}Checking Kubernetes secrets...${NC}"
    if kubectl get secret cloudflare-secrets -n naytife -o yaml | grep -q "templates-bucket-name"; then
        echo -e "${GREEN}✅ Templates bucket name configured in secrets${NC}"
    else
        echo -e "${RED}❌ Templates bucket name not found in secrets${NC}"
    fi
    
    if kubectl get secret cloudflare-secrets -n naytife -o yaml | grep -q "stores-bucket-name"; then
        echo -e "${GREEN}✅ Stores bucket name configured in secrets${NC}"
    else
        echo -e "${RED}❌ Stores bucket name not found in secrets${NC}"
    fi
    
    # Decode and show bucket names
    templates_bucket=$(kubectl get secret cloudflare-secrets -n naytife -o jsonpath='{.data.templates-bucket-name}' | base64 -d 2>/dev/null || echo "not found")
    stores_bucket=$(kubectl get secret cloudflare-secrets -n naytife -o jsonpath='{.data.stores-bucket-name}' | base64 -d 2>/dev/null || echo "not found")
    
    echo -e "   Templates bucket: ${BLUE}$templates_bucket${NC}"
    echo -e "   Stores bucket: ${BLUE}$stores_bucket${NC}"
}

# Show recommendations
show_recommendations() {
    echo -e "\n${YELLOW}💡 Recommendations${NC}"
    echo -e "${YELLOW}=================${NC}"
    echo -e "1. ${BLUE}Monitor service logs after deployment:${NC}"
    echo -e "   kubectl logs deployment/template-registry -n naytife -f"
    echo -e "   kubectl logs deployment/store-deployer -n naytife -f"
    echo -e ""
    echo -e "2. ${BLUE}Test functionality:${NC}"
    echo -e "   - Try uploading a template"
    echo -e "   - Try deploying a store"
    echo -e "   - Upload product/shop images via dashboard"
    echo -e ""
    echo -e "3. ${BLUE}Update dashboard environment variables:${NC}"
    echo -e "   CLOUDFLARE_R2_BUCKET=stores"
    echo -e ""
    echo -e "4. ${BLUE}Monitor bucket usage in Cloudflare R2 dashboard${NC}"
}

# Main execution
main() {
    echo -e "${BLUE}Running bucket separation tests...${NC}\n"
    
    check_environment
    test_kubernetes_services
    test_api_endpoints
    test_bucket_access
    show_recommendations
    
    echo -e "\n${GREEN}🎉 Testing completed!${NC}"
    echo -e "${BLUE}Review the results above and address any issues found.${NC}"
}

# Run the main function
main "$@"
