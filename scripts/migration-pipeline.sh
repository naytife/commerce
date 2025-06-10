#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔄 Database Migration CI/CD Pipeline${NC}"
echo "======================================"

# Configuration
NAMESPACE="naytife"
BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../backend" && pwd)"
MIGRATIONS_DIR="$BACKEND_DIR/internal/db/migrations"
MAX_WAIT_TIME=600  # 10 minutes
CHECK_INTERVAL=10  # 10 seconds

# Function to check if migrations exist
check_migrations_exist() {
    echo -e "${YELLOW}📁 Checking for migration files...${NC}"
    
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        echo -e "${RED}❌ Migration directory not found: $MIGRATIONS_DIR${NC}"
        exit 1
    fi
    
    MIGRATION_COUNT=$(find "$MIGRATIONS_DIR" -name "*.sql" | wc -l)
    if [ "$MIGRATION_COUNT" -eq 0 ]; then
        echo -e "${YELLOW}⚠️  No migration files found, skipping migration step${NC}"
        exit 0
    fi
    
    echo -e "${GREEN}✅ Found $MIGRATION_COUNT migration files${NC}"
    ls -la "$MIGRATIONS_DIR"/*.sql
}

# Function to validate migration files
validate_migrations() {
    echo -e "${YELLOW}🔍 Validating migration files...${NC}"
    
    cd "$BACKEND_DIR"
    
    # Check if atlas is available
    if ! command -v atlas &> /dev/null; then
        echo -e "${YELLOW}⚠️  Atlas not found, installing...${NC}"
        # Install atlas if not available
        curl -sSf https://atlasgo.sh | sh
        export PATH="$HOME/.atlas:$PATH"
    fi
    
    # Validate migrations
    if atlas migrate validate --env local; then
        echo -e "${GREEN}✅ Migration files are valid${NC}"
    else
        echo -e "${RED}❌ Migration validation failed${NC}"
        exit 1
    fi
    
    # Check for migration hash consistency
    if atlas migrate hash --env local; then
        echo -e "${GREEN}✅ Migration hashes are consistent${NC}"
    else
        echo -e "${RED}❌ Migration hash check failed${NC}"
        exit 1
    fi
}

# Function to update ConfigMap
update_migration_configmap() {
    echo -e "${YELLOW}📝 Updating migration ConfigMap...${NC}"
    
    # Delete existing ConfigMap if it exists
    kubectl delete configmap backend-migrations -n $NAMESPACE 2>/dev/null || true
    
    # Create new ConfigMap with migration files
    kubectl create configmap backend-migrations -n $NAMESPACE \
        --from-file="$MIGRATIONS_DIR" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ Migration ConfigMap updated${NC}"
}

# Function to run migration job
run_migration_job() {
    echo -e "${YELLOW}🚀 Starting migration job...${NC}"
    
    # Clean up any existing migration job
    kubectl delete job backend-migrate -n $NAMESPACE 2>/dev/null || true
    
    # Wait a moment for cleanup
    sleep 5
    
    # Apply migration job
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    MANIFESTS_DIR="$(cd "$SCRIPT_DIR/../k3s/manifests" && pwd)"
    kubectl apply -f "$MANIFESTS_DIR/06-backend/backend-migration.yaml"
    
    echo -e "${GREEN}✅ Migration job started${NC}"
}

# Function to monitor migration job
monitor_migration_job() {
    echo -e "${YELLOW}⏳ Monitoring migration job progress...${NC}"
    
    local start_time=$(date +%s)
    local current_time
    local elapsed_time
    
    while true; do
        current_time=$(date +%s)
        elapsed_time=$((current_time - start_time))
        
        # Check if max wait time exceeded
        if [ $elapsed_time -gt $MAX_WAIT_TIME ]; then
            echo -e "${RED}❌ Migration job timed out after ${MAX_WAIT_TIME} seconds${NC}"
            kubectl logs job/backend-migrate -n $NAMESPACE --tail=50
            exit 1
        fi
        
        # Check job status
        if kubectl get job backend-migrate -n $NAMESPACE &> /dev/null; then
            JOB_COMPLETE=$(kubectl get job backend-migrate -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' 2>/dev/null || echo "False")
            JOB_FAILED=$(kubectl get job backend-migrate -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Failed")].status}' 2>/dev/null || echo "False")
            
            if [[ "$JOB_COMPLETE" == "True" ]]; then
                echo -e "${GREEN}✅ Migration job completed successfully${NC}"
                kubectl logs job/backend-migrate -n $NAMESPACE --tail=20
                return 0
            elif [[ "$JOB_FAILED" == "True" ]]; then
                echo -e "${RED}❌ Migration job failed${NC}"
                kubectl logs job/backend-migrate -n $NAMESPACE --tail=50
                kubectl describe job backend-migrate -n $NAMESPACE
                exit 1
            else
                echo -e "${BLUE}⏳ Migration job still running... (${elapsed_time}s elapsed)${NC}"
                # Show recent logs
                kubectl logs job/backend-migrate -n $NAMESPACE --tail=5 2>/dev/null || true
            fi
        else
            echo -e "${YELLOW}⚠️  Migration job not found, waiting...${NC}"
        fi
        
        sleep $CHECK_INTERVAL
    done
}

# Function to verify migration success
verify_migration_success() {
    echo -e "${YELLOW}🔍 Verifying migration success...${NC}"
    
    # Check if we can connect to the database
    if kubectl run db-check --rm -i --restart=Never --image=postgres:15-alpine -n $NAMESPACE -- \
        sh -c "PGPASSWORD=\$POSTGRES_PASSWORD psql -h postgres.naytife.svc.cluster.local -U postgres -d naytifedb -c 'SELECT version();'" \
        --env POSTGRES_PASSWORD="$(kubectl get secret postgres-secret -n $NAMESPACE -o jsonpath='{.data.postgres-password}' | base64 -d)" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Database connection verified${NC}"
    else
        echo -e "${RED}❌ Database connection failed${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Migration pipeline completed successfully${NC}"
}

# Main pipeline execution
main() {
    echo -e "${BLUE}Starting migration pipeline...${NC}"
    
    # Step 1: Check prerequisites
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl is required but not installed${NC}"
        exit 1
    fi
    
    if ! kubectl get namespace $NAMESPACE &> /dev/null; then
        echo -e "${RED}❌ Namespace $NAMESPACE not found${NC}"
        exit 1
    fi
    
    # Step 2: Check and validate migrations
    check_migrations_exist
    validate_migrations
    
    # Step 3: Update ConfigMap
    update_migration_configmap
    
    # Step 4: Run migration job
    run_migration_job
    
    # Step 5: Monitor migration job
    monitor_migration_job
    
    # Step 6: Verify success
    verify_migration_success
    
    echo -e "${GREEN}🎉 Migration pipeline completed successfully!${NC}"
}

# Run the pipeline
main "$@"
