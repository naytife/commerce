#!/bin/bash

# Bucket Separation Migration Script
# This script helps migrate from single naytife-shops-static bucket to separate templates and stores buckets

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
OLD_BUCKET="naytife-shops-static"
TEMPLATES_BUCKET="templates"
STORES_BUCKET="stores"

echo -e "${BLUE}🚀 Naytife Bucket Separation Migration${NC}"
echo -e "${BLUE}=====================================${NC}"

# Function to check if rclone is installed
check_rclone() {
    if ! command -v rclone &> /dev/null; then
        echo -e "${RED}❌ rclone is not installed. Please install it first:${NC}"
        echo "   brew install rclone"
        echo "   or visit: https://rclone.org/install/"
        exit 1
    fi
    echo -e "${GREEN}✅ rclone is available${NC}"
}

# Function to check if buckets exist
check_buckets() {
    echo -e "\n${YELLOW}📦 Checking bucket availability...${NC}"
    
    # Check old bucket
    if rclone lsd r2: | grep -q "$OLD_BUCKET"; then
        echo -e "${GREEN}✅ Found source bucket: $OLD_BUCKET${NC}"
    else
        echo -e "${RED}❌ Source bucket $OLD_BUCKET not found${NC}"
        exit 1
    fi
    
    # Check if new buckets exist
    for bucket in "$TEMPLATES_BUCKET" "$STORES_BUCKET"; do
        if rclone lsd r2: | grep -q "$bucket"; then
            echo -e "${GREEN}✅ Target bucket exists: $bucket${NC}"
        else
            echo -e "${YELLOW}⚠️  Target bucket does not exist: $bucket${NC}"
            echo -e "${YELLOW}   Please create it in your Cloudflare R2 dashboard first${NC}"
            exit 1
        fi
    done
}

# Function to show migration preview
show_migration_preview() {
    echo -e "\n${YELLOW}🔍 Migration Preview${NC}"
    echo -e "${YELLOW}===================${NC}"
    
    echo -e "\n${BLUE}Templates to migrate:${NC}"
    if rclone ls "r2:$OLD_BUCKET/templates/" 2>/dev/null | head -10; then
        echo "   ... (showing first 10 files)"
    else
        echo "   No templates found"
    fi
    
    echo -e "\n${BLUE}Stores to migrate:${NC}"
    if rclone ls "r2:$OLD_BUCKET/stores/" 2>/dev/null | head -10; then
        echo "   ... (showing first 10 files)"
    else
        echo "   No stores found"
    fi
    
    echo -e "\n${BLUE}Other files (products/shops):${NC}"
    rclone ls "r2:$OLD_BUCKET/" 2>/dev/null | grep -E "(^.*products/|^.*shops/)" | head -10 || echo "   No product/shop files found"
}

# Function to migrate templates
migrate_templates() {
    echo -e "\n${YELLOW}📁 Migrating templates...${NC}"
    
    if rclone ls "r2:$OLD_BUCKET/templates/" &>/dev/null; then
        echo -e "${BLUE}Copying templates from $OLD_BUCKET/templates/ to $TEMPLATES_BUCKET/${NC}"
        rclone copy "r2:$OLD_BUCKET/templates/" "r2:$TEMPLATES_BUCKET/" --progress --stats 30s
        
        # Verify
        template_count=$(rclone ls "r2:$TEMPLATES_BUCKET/" | wc -l)
        echo -e "${GREEN}✅ Migrated $template_count template files${NC}"
    else
        echo -e "${YELLOW}⚠️  No templates found to migrate${NC}"
    fi
}

# Function to migrate stores
migrate_stores() {
    echo -e "\n${YELLOW}🏪 Migrating stores...${NC}"
    
    # Migrate store files
    if rclone ls "r2:$OLD_BUCKET/stores/" &>/dev/null; then
        echo -e "${BLUE}Copying stores from $OLD_BUCKET/stores/ to $STORES_BUCKET/${NC}"
        rclone copy "r2:$OLD_BUCKET/stores/" "r2:$STORES_BUCKET/" --progress --stats 30s
        
        store_count=$(rclone ls "r2:$STORES_BUCKET/" | wc -l)
        echo -e "${GREEN}✅ Migrated $store_count store files${NC}"
    else
        echo -e "${YELLOW}⚠️  No stores found to migrate${NC}"
    fi
    
    # Migrate product images
    if rclone ls "r2:$OLD_BUCKET/products/" &>/dev/null; then
        echo -e "${BLUE}Copying products from $OLD_BUCKET/products/ to $STORES_BUCKET/products/${NC}"
        rclone copy "r2:$OLD_BUCKET/products/" "r2:$STORES_BUCKET/products/" --progress --stats 30s
        
        product_count=$(rclone ls "r2:$STORES_BUCKET/products/" | wc -l)
        echo -e "${GREEN}✅ Migrated $product_count product files${NC}"
    else
        echo -e "${YELLOW}⚠️  No product files found to migrate${NC}"
    fi
    
    # Migrate shop images
    if rclone ls "r2:$OLD_BUCKET/shops/" &>/dev/null; then
        echo -e "${BLUE}Copying shops from $OLD_BUCKET/shops/ to $STORES_BUCKET/shops/${NC}"
        rclone copy "r2:$OLD_BUCKET/shops/" "r2:$STORES_BUCKET/shops/" --progress --stats 30s
        
        shop_count=$(rclone ls "r2:$STORES_BUCKET/shops/" | wc -l)
        echo -e "${GREEN}✅ Migrated $shop_count shop files${NC}"
    else
        echo -e "${YELLOW}⚠️  No shop files found to migrate${NC}"
    fi
}

# Function to update Kubernetes services
update_kubernetes_services() {
    echo -e "\n${YELLOW}☸️  Updating Kubernetes services...${NC}"
    
    # Apply the updated manifests
    echo -e "${BLUE}Applying updated secrets and manifests...${NC}"
    kubectl apply -f k3s/manifests/08-template-system/cloudflare-secrets.yaml
    kubectl apply -f k3s/manifests/08-template-system/template-registry.yaml
    kubectl apply -f k3s/manifests/08-template-system/store-deployer.yaml
    
    # Restart services to pick up new environment variables
    echo -e "${BLUE}Restarting services...${NC}"
    kubectl rollout restart deployment/template-registry -n naytife
    kubectl rollout restart deployment/store-deployer -n naytife
    
    # Wait for deployments to be ready
    echo -e "${BLUE}Waiting for deployments to be ready...${NC}"
    kubectl rollout status deployment/template-registry -n naytife --timeout=300s
    kubectl rollout status deployment/store-deployer -n naytife --timeout=300s
    
    echo -e "${GREEN}✅ Kubernetes services updated${NC}"
}

# Function to verify migration
verify_migration() {
    echo -e "\n${YELLOW}🔍 Verifying migration...${NC}"
    
    # Check service logs
    echo -e "${BLUE}Checking template-registry logs...${NC}"
    if kubectl logs deployment/template-registry -n naytife --tail=5 | grep -i "bucket.*templates"; then
        echo -e "${GREEN}✅ Template-registry is using templates bucket${NC}"
    else
        echo -e "${RED}❌ Template-registry bucket configuration unclear${NC}"
    fi
    
    echo -e "${BLUE}Checking store-deployer logs...${NC}"
    if kubectl logs deployment/store-deployer -n naytife --tail=5 | grep -i "bucket.*stores"; then
        echo -e "${GREEN}✅ Store-deployer is using stores bucket${NC}"
    else
        echo -e "${RED}❌ Store-deployer bucket configuration unclear${NC}"
    fi
    
    # Check bucket contents
    templates_count=$(rclone ls "r2:$TEMPLATES_BUCKET/" | wc -l)
    stores_count=$(rclone ls "r2:$STORES_BUCKET/" | wc -l)
    
    echo -e "\n${BLUE}📊 Migration Summary:${NC}"
    echo -e "   Templates bucket: $templates_count files"
    echo -e "   Stores bucket: $stores_count files"
}

# Function to show next steps
show_next_steps() {
    echo -e "\n${YELLOW}📝 Next Steps${NC}"
    echo -e "${YELLOW}============${NC}"
    echo -e "1. ${BLUE}Update your dashboard environment variables:${NC}"
    echo -e "   CLOUDFLARE_R2_BUCKET=stores"
    echo -e ""
    echo -e "2. ${BLUE}Test the migration:${NC}"
    echo -e "   - Upload a new template"
    echo -e "   - Deploy a store"
    echo -e "   - Upload product/shop images"
    echo -e ""
    echo -e "3. ${BLUE}Monitor service logs:${NC}"
    echo -e "   kubectl logs deployment/template-registry -n naytife -f"
    echo -e "   kubectl logs deployment/store-deployer -n naytife -f"
    echo -e ""
    echo -e "4. ${BLUE}After verification, clean up old bucket (optional):${NC}"
    echo -e "   rclone delete r2:$OLD_BUCKET/templates/ --dry-run"
    echo -e "   rclone delete r2:$OLD_BUCKET/stores/ --dry-run"
    echo -e "   # Remove --dry-run when ready to delete"
}

# Main execution
main() {
    echo -e "${BLUE}Starting bucket separation migration...${NC}"
    
    # Parse command line arguments
    DRY_RUN=false
    SKIP_KUBERNETES=false
    
    for arg in "$@"; do
        case $arg in
            --dry-run)
                DRY_RUN=true
                echo -e "${YELLOW}🔍 DRY RUN MODE - No actual changes will be made${NC}"
                ;;
            --skip-kubernetes)
                SKIP_KUBERNETES=true
                echo -e "${YELLOW}⏭️  SKIPPING Kubernetes updates${NC}"
                ;;
            --help)
                echo "Usage: $0 [--dry-run] [--skip-kubernetes] [--help]"
                echo ""
                echo "Options:"
                echo "  --dry-run          Show what would be done without making changes"
                echo "  --skip-kubernetes  Skip Kubernetes service updates"
                echo "  --help            Show this help message"
                exit 0
                ;;
        esac
    done
    
    # Pre-flight checks
    check_rclone
    check_buckets
    show_migration_preview
    
    # Confirm before proceeding
    if [ "$DRY_RUN" = false ]; then
        echo -e "\n${YELLOW}⚠️  This will migrate data from $OLD_BUCKET to separate buckets.${NC}"
        read -p "Do you want to proceed? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${RED}❌ Migration cancelled${NC}"
            exit 1
        fi
    fi
    
    if [ "$DRY_RUN" = false ]; then
        # Perform migration
        migrate_templates
        migrate_stores
        
        if [ "$SKIP_KUBERNETES" = false ]; then
            update_kubernetes_services
        fi
        
        verify_migration
        show_next_steps
        
        echo -e "\n${GREEN}🎉 Migration completed successfully!${NC}"
    else
        echo -e "\n${BLUE}🔍 DRY RUN completed - no changes made${NC}"
        echo -e "${BLUE}Run without --dry-run to perform the actual migration${NC}"
    fi
}

# Run the main function with all arguments
main "$@"
