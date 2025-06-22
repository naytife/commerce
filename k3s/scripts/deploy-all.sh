#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Deploying Naytife Services to k3s${NC}"
echo "====================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFESTS_DIR="$(cd "$SCRIPT_DIR/../manifests" && pwd)"

# Function to wait for deployment to be ready
wait_for_deployment() {
    local namespace=$1
    local deployment=$2
    local timeout=${3:-300}
    
    echo -e "${YELLOW}⏳ Waiting for $deployment in namespace $namespace to be ready...${NC}"
    
    if kubectl wait --for=condition=available deployment/$deployment -n $namespace --timeout=${timeout}s >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $deployment is ready${NC}"
        return 0
    else
        echo -e "${RED}❌ $deployment failed to become ready within ${timeout}s${NC}"
        return 1
    fi
}

# Function to wait for pods to be ready
wait_for_pods() {
    local namespace=$1
    local label_selector=$2
    local timeout=${3:-300}
    
    echo -e "${YELLOW}⏳ Waiting for pods with label $label_selector in namespace $namespace...${NC}"
    
    if kubectl wait --for=condition=ready pod -l $label_selector -n $namespace --timeout=${timeout}s >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Pods are ready${NC}"
        return 0
    else
        echo -e "${RED}❌ Pods failed to become ready within ${timeout}s${NC}"
        return 1
    fi
}

# Function to wait for job completion
wait_for_job() {
    local namespace=$1
    local job_name=$2
    local timeout=${3:-300}
    
    echo -e "${YELLOW}⏳ Waiting for job $job_name in namespace $namespace to complete...${NC}"
    
    if kubectl wait --for=condition=complete job/$job_name -n $namespace --timeout=${timeout}s >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Job $job_name completed successfully${NC}"
        return 0
    else
        echo -e "${RED}❌ Job $job_name failed to complete within ${timeout}s${NC}"
        echo "Job logs:"
        kubectl logs job/$job_name -n $namespace
        return 1
    fi
}

# Deploy services in order
echo -e "\n${BLUE}📦 Step 1: Creating Namespaces${NC}"
kubectl apply -f "$MANIFESTS_DIR/00-namespaces/"
echo -e "${GREEN}✅ Namespaces created${NC}"

echo -e "\n${BLUE}🐘 Step 2: Deploying PostgreSQL${NC}"
kubectl apply -f "$MANIFESTS_DIR/01-postgres/"
wait_for_deployment "naytife" "postgres" 200

echo -e "\n${BLUE}📊 Step 3: Deploying Redis${NC}"
kubectl apply -f "$MANIFESTS_DIR/02-redis/"
wait_for_deployment "naytife" "redis" 120

echo -e "\n${BLUE}🆔 Step 4: Deploying Hydra${NC}"
kubectl apply -f "$MANIFESTS_DIR/03-hydra/"
# Wait for migration to complete first
echo -e "${YELLOW}⏳ Waiting for Hydra migration to complete...${NC}"
sleep 10
wait_for_deployment "naytife-auth" "hydra" 180

echo -e "\n${BLUE}🔐 Step 5: Deploying Oathkeeper${NC}"
kubectl apply -f "$MANIFESTS_DIR/04-oathkeeper/"
wait_for_deployment "naytife-auth" "oathkeeper" 120

echo -e "\n${BLUE}🔑 Step 6: Deploying Auth Handler${NC}"
kubectl apply -f "$MANIFESTS_DIR/05-auth-handler/"
wait_for_deployment "naytife-auth" "auth-handler" 120

echo -e "\n${BLUE}🔙 Step 7: Deploying Backend API${NC}"

# Deploy backend configuration and secrets first (needed for migrations)
echo -e "${YELLOW}🔐 Applying backend configuration...${NC}"
kubectl apply -f "$MANIFESTS_DIR/06-backend/backend.yaml"
echo -e "${GREEN}✅ Backend configuration applied${NC}"

# First, create migration ConfigMap with actual migration files
echo -e "${YELLOW}📝 Creating migration ConfigMap...${NC}"
BACKEND_DIR="$(cd "$SCRIPT_DIR/../../backend" && pwd)"
MIGRATIONS_DIR="$BACKEND_DIR/internal/db/migrations"

if [ -d "$MIGRATIONS_DIR" ]; then
    # Delete existing migration ConfigMap if it exists
    kubectl delete configmap backend-migrations -n naytife 2>/dev/null || true
    
    # Create ConfigMap with migration files
    kubectl create configmap backend-migrations -n naytife \
        --from-file="$MIGRATIONS_DIR" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ Migration ConfigMap created${NC}"
else
    echo -e "${RED}❌ Migration directory not found at $MIGRATIONS_DIR${NC}"
    exit 1
fi

# Now run the migration job
echo -e "${YELLOW}🗄️  Running database migrations...${NC}"
kubectl apply -f "$MANIFESTS_DIR/06-backend/backend-migration.yaml"

# Wait for migration job to complete
echo -e "${YELLOW}⏳ Waiting for migration job to complete...${NC}"
kubectl wait --for=condition=complete job/backend-migrate -n naytife --timeout=300s

# Check if migration was successful
if kubectl get job backend-migrate -n naytife -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' | grep -q "True"; then
    echo -e "${GREEN}✅ Database migration completed successfully${NC}"
else
    echo -e "${RED}❌ Database migration failed${NC}"
    echo "Migration job logs:"
    kubectl logs job/backend-migrate -n naytife
    exit 1
fi

# Wait for backend deployment to be ready
wait_for_deployment "naytife" "backend" 180

echo -e "\n${BLUE}🏗️  Step 8: Deploying Template System${NC}"
kubectl apply -f "$MANIFESTS_DIR/08-template-system/"

# Wait for template system deployments
echo -e "${YELLOW}⏳ Waiting for template system deployments to be ready...${NC}"
wait_for_deployment "naytife" "template-registry" 180
wait_for_deployment "naytife" "store-deployer" 180

echo -e "${GREEN}✅ Template System deployed successfully${NC}"

echo -e "\n${GREEN}🎉 All services deployed successfully!${NC}"

echo -e "\n${BLUE}📊 Deployment Status:${NC}"
echo "================================"
kubectl get pods --all-namespaces -l app.kubernetes.io/part-of=naytife-platform

echo -e "\n${BLUE}🔗 Service Access Points:${NC}"
echo "================================"
echo "  🔐 API Gateway:      http://127.0.0.1:8080"
echo "  🔙 Backend API:      http://127.0.0.1:8000"
echo "  🔑 Auth Handler:     http://127.0.0.1:3000"
echo "  🏗️  Template Registry: http://127.0.0.1:9001"
echo "  🚀 Store Deployer:   http://127.0.0.1:9003"
echo "  🐘 PostgreSQL:       localhost:5432"
echo "  📊 Redis:            localhost:6379"
echo "  🆔 Hydra Public:     http://127.0.0.1:4444"
echo "  🆔 Hydra Admin:      http://127.0.0.1:4445"

echo -e "\n${BLUE}📋 Quick Health Check:${NC}"
echo "================================"

# Wait a moment for services to be fully ready
sleep 5

echo -n "🔐 Oathkeeper:   "
if curl -s http://127.0.0.1:8080/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

echo -n "🔙 Backend:      "
if curl -s http://127.0.0.1:8000/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

echo -n "🔑 Auth Handler:     "
if curl -s http://127.0.0.1:3000/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

echo -n "🏗️  Template Registry: "
if curl -s http://127.0.0.1:9001/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

echo -n "🚀 Store Deployer:   "
if curl -s http://127.0.0.1:9003/health >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Healthy${NC}"
else
    echo -e "${RED}❌ Not ready${NC}"
fi

echo -e "\n${BLUE}📝 Next Steps:${NC}"
echo "  • Test the deployment: ./scripts/test-deployment.sh"
echo "  • View logs: ./scripts/logs.sh [service-name]"
echo "  • Check status: ./scripts/status.sh"
echo "  • API Documentation: http://127.0.0.1:8080/v1/docs"
