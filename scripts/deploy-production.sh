#!/bin/bash

# Production Deployment Script for Naytife Commerce Platform
# This script deploys all services to a production k3s cluster
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Function to load environment variables from .env file
load_env() {
    if [[ -f .env ]]; then
        echo -e "${GREEN}Loading environment variables from .env file...${NC}"
        export $(grep -v '^#' .env | xargs)
    else
        echo -e "${YELLOW}Warning: .env file not found for production deployment.${NC}"
        echo -e "${RED}Production deployment requires properly configured environment variables.${NC}"
        exit 1
    fi
}

echo -e "${BLUE}🚀 Naytife Production Deployment${NC}"
echo "=================================="

# Load environment variables
load_env

# Production Configuration (with .env fallbacks)
CLUSTER_NAME="naytife-prod"
DOMAIN="${DOMAIN:-yourdomain.com}"
DB_PASSWORD="${POSTGRES_PASSWORD:-$(openssl rand -base64 32)}"
REDIS_PASSWORD="${REDIS_PASSWORD:-$(openssl rand -base64 32)}"
HYDRA_SYSTEM_SECRET="${HYDRA_SYSTEM_SECRET:-$(openssl rand -base64 32)}"
GOOGLE_CLIENT_ID="${GOOGLE_CLIENT_ID}"
GOOGLE_CLIENT_SECRET="${GOOGLE_CLIENT_SECRET}"
STRIPE_SECRET_KEY="${STRIPE_SECRET_KEY}"
CLOUDFLARE_R2_ACCESS_KEY="${R2_ACCESS_KEY}"
CLOUDFLARE_R2_SECRET_KEY="${R2_SECRET_KEY}"

# Validate required environment variables
validate_env() {
    local missing=()
    
    [ -z "$DOMAIN" ] && missing+=("DOMAIN")
    [ -z "$GOOGLE_CLIENT_ID" ] && missing+=("GOOGLE_CLIENT_ID")
    [ -z "$GOOGLE_CLIENT_SECRET" ] && missing+=("GOOGLE_CLIENT_SECRET")
    [ -z "$STRIPE_SECRET_KEY" ] && missing+=("STRIPE_SECRET_KEY")
    [ -z "$CLOUDFLARE_R2_ACCESS_KEY" ] && missing+=("CLOUDFLARE_R2_ACCESS_KEY")
    [ -z "$CLOUDFLARE_R2_SECRET_KEY" ] && missing+=("CLOUDFLARE_R2_SECRET_KEY")
    
    if [ ${#missing[@]} -ne 0 ]; then
        echo -e "${RED}❌ Missing required environment variables:${NC}"
        printf '%s\n' "${missing[@]}"
        echo ""
        echo "Please set these variables before running the script:"
        echo "export DOMAIN=yourdomain.com"
        echo "export GOOGLE_CLIENT_ID=your_google_client_id"
        echo "export GOOGLE_CLIENT_SECRET=your_google_client_secret"
        echo "export STRIPE_SECRET_KEY=sk_live_your_stripe_secret"
        echo "export CLOUDFLARE_R2_ACCESS_KEY=your_r2_access_key"
        echo "export CLOUDFLARE_R2_SECRET_KEY=your_r2_secret_key"
        exit 1
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to wait for deployment to be ready
wait_for_deployment() {
    local namespace=$1
    local deployment=$2
    local timeout=600

    echo -e "${YELLOW}⏳ Waiting for $deployment in $namespace to be ready...${NC}"
    
    if kubectl wait --for=condition=available deployment/$deployment -n $namespace --timeout=${timeout}s; then
        echo -e "${GREEN}✅ $deployment is ready${NC}"
        return 0
    else
        echo -e "${RED}❌ $deployment failed to be ready within ${timeout}s${NC}"
        kubectl describe deployment/$deployment -n $namespace
        return 1
    fi
}

# Step 1: Validation
echo -e "\n${BLUE}Step 1: Environment Validation${NC}"
echo "----------------------------------------"

validate_env

echo -e "${GREEN}✅ Environment variables validated${NC}"
echo "Domain: $DOMAIN"
echo "Database Password: ${DB_PASSWORD:0:8}..."
echo "Redis Password: ${REDIS_PASSWORD:0:8}..."

# Step 2: Prerequisites Check
echo -e "\n${BLUE}Step 2: Checking Prerequisites${NC}"
echo "----------------------------------------"

MISSING_DEPS=()
for cmd in k3s kubectl helm docker; do
    if command_exists $cmd; then
        echo -e "${GREEN}✅ $cmd is installed${NC}"
    else
        echo -e "${RED}❌ $cmd is not installed${NC}"
        MISSING_DEPS+=($cmd)
    fi
done

if [ ${#MISSING_DEPS[@]} -ne 0 ]; then
    echo -e "\n${RED}Missing dependencies: ${MISSING_DEPS[*]}${NC}"
    echo "Please install missing dependencies and run this script again."
    exit 1
fi

# Step 3: Setup k3s (if not already running)
echo -e "\n${BLUE}Step 3: Setting up k3s${NC}"
echo "----------------------------------------"

if systemctl is-active --quiet k3s; then
    echo -e "${GREEN}✅ k3s is already running${NC}"
else
    echo "🔧 Installing and starting k3s..."
    curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable traefik" sh -
    
    # Wait for k3s to be ready
    sleep 30
    
    # Setup kubeconfig for current user
    mkdir -p ~/.kube
    sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
    sudo chown $(id -u):$(id -g) ~/.kube/config
    
    echo -e "${GREEN}✅ k3s installed and configured${NC}"
fi

# Step 4: Setup Helm
echo -e "\n${BLUE}Step 4: Setting up Helm Repositories${NC}"
echo "----------------------------------------"

if ! command_exists helm; then
    echo "📦 Installing Helm..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
fi

echo "📈 Adding Helm repositories..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add cert-manager https://charts.jetstack.io
helm repo update

echo -e "${GREEN}✅ Helm repositories configured${NC}"

# Step 5: Install Ingress Controller
echo -e "\n${BLUE}Step 5: Installing Ingress Controller${NC}"
echo "----------------------------------------"

if kubectl get namespace ingress-nginx >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Ingress controller already installed${NC}"
else
    echo "🌐 Installing NGINX Ingress Controller..."
    helm install ingress-nginx ingress-nginx/ingress-nginx \
        --create-namespace \
        --namespace ingress-nginx \
        --set controller.service.type=LoadBalancer \
        --wait
    
    echo -e "${GREEN}✅ Ingress controller installed${NC}"
fi

# Step 6: Install Cert-Manager for SSL
echo -e "\n${BLUE}Step 6: Installing Cert-Manager${NC}"
echo "----------------------------------------"

if kubectl get namespace cert-manager >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Cert-manager already installed${NC}"
else
    echo "🔒 Installing Cert-Manager..."
    helm install cert-manager cert-manager/cert-manager \
        --create-namespace \
        --namespace cert-manager \
        --set installCRDs=true \
        --wait
    
    echo -e "${GREEN}✅ Cert-manager installed${NC}"
fi

# Step 7: Create Namespaces and Secrets
echo -e "\n${BLUE}Step 7: Creating Namespaces and Secrets${NC}"
echo "----------------------------------------"

# Create namespaces
for ns in database auth backend cloud-build; do
    if kubectl get namespace $ns >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Namespace $ns already exists${NC}"
    else
        kubectl create namespace $ns
        echo -e "${GREEN}✅ Created namespace $ns${NC}"
    fi
done

# Create production secrets
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: auth-handler-secrets
  namespace: auth
type: Opaque
stringData:
  HYDRA_ADMIN_URL: "http://hydra-admin.auth.svc.cluster.local:4445"
  GOOGLE_CLIENT_ID: "$GOOGLE_CLIENT_ID"
  GOOGLE_CLIENT_SECRET: "$GOOGLE_CLIENT_SECRET"
  LOGIN_HANDLER_REDIRECT_URI: "https://$DOMAIN/callback"
---
apiVersion: v1
kind: Secret
metadata:
  name: backend-secrets
  namespace: backend
type: Opaque
stringData:
  DATABASE_URL: "postgres://postgres:$DB_PASSWORD@postgres-postgresql.database.svc.cluster.local:5432/naytifedb?sslmode=disable"
  REDIS_URL: "redis://:$REDIS_PASSWORD@redis-master.database.svc.cluster.local:6379"
  STRIPE_SECRET_KEY: "$STRIPE_SECRET_KEY"
  JWT_SECRET: "$HYDRA_SYSTEM_SECRET"
---
apiVersion: v1
kind: Secret
metadata:
  name: cloud-build-secrets
  namespace: cloud-build
type: Opaque
stringData:
  REDIS_URL: "redis://:$REDIS_PASSWORD@redis-master.database.svc.cluster.local:6379"
  CLOUDFLARE_R2_ACCESS_KEY: "$CLOUDFLARE_R2_ACCESS_KEY"
  CLOUDFLARE_R2_SECRET_KEY: "$CLOUDFLARE_R2_SECRET_KEY"
  CLOUDFLARE_R2_BUCKET: "naytife-sites"
  CLOUDFLARE_R2_ENDPOINT: "https://your-account-id.r2.cloudflarestorage.com"
EOF

echo -e "${GREEN}✅ Secrets created${NC}"

# Step 8: Deploy Infrastructure Services
echo -e "\n${BLUE}Step 8: Deploying Infrastructure Services${NC}"
echo "----------------------------------------"

# Deploy PostgreSQL
if kubectl get deployment postgres-postgresql -n database >/dev/null 2>&1; then
    echo -e "${GREEN}✅ PostgreSQL already deployed${NC}"
else
    echo "🐘 Deploying PostgreSQL..."
    helm install postgres bitnami/postgresql \
        --namespace database \
        --set auth.postgresPassword=$DB_PASSWORD \
        --set auth.database=naytifedb \
        --set primary.persistence.size=20Gi \
        --set primary.persistence.storageClass=local-path \
        --wait
    
    echo -e "${GREEN}✅ PostgreSQL deployed${NC}"
fi

# Deploy Redis
if kubectl get deployment redis-master -n database >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Redis already deployed${NC}"
else
    echo "📊 Deploying Redis..."
    helm install redis bitnami/redis \
        --namespace database \
        --set auth.password=$REDIS_PASSWORD \
        --set master.persistence.size=5Gi \
        --set master.persistence.storageClass=local-path \
        --wait
    
    echo -e "${GREEN}✅ Redis deployed${NC}"
fi

# Step 9: Deploy Custom Services
echo -e "\n${BLUE}Step 9: Deploying Custom Services${NC}"
echo "----------------------------------------"

# Build and deploy Hydra
echo "🔐 Deploying Ory Hydra..."
# Update Hydra config for production
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: hydra-config
  namespace: auth
data:
  hydra.yaml: |
    log:
      level: info
    
    dsn: postgres://postgres:$DB_PASSWORD@postgres-postgresql.database.svc.cluster.local:5432/naytifedb?search_path=hydra&sslmode=disable
    
    urls:
      self:
        issuer: https://$DOMAIN/
        admin: http://hydra-admin.auth.svc.cluster.local:4445
      login: https://$DOMAIN/login
      logout: https://$DOMAIN/logout
      consent: https://$DOMAIN/consent
      post_logout_redirect: https://$DOMAIN/
    
    serve:
      public:
        port: 4444
        host: 0.0.0.0
        cors:
          enabled: true
          allowed_origins:
            - https://$DOMAIN
          allowed_methods:
            - POST
            - GET
            - PUT
            - PATCH
            - DELETE
            - OPTIONS
          allowed_headers:
            - Authorization
            - Content-Type
            - X-Requested-With
          exposed_headers:
            - Content-Type
          allow_credentials: true
      admin:
        port: 4445
        host: 0.0.0.0
      cookies:
        domain: .$DOMAIN
        same_site_mode: Lax
        secure: true
    
    secrets:
      system:
        - $HYDRA_SYSTEM_SECRET
    
    oidc:
      subject_identifiers:
        supported_types:
          - pairwise
          - public
        pairwise:
          salt: $HYDRA_SYSTEM_SECRET
    
    oauth2:
      expose_internal_errors: false
      hashers:
        algorithm: bcrypt
        bcrypt:
          cost: 12
    
    strategies:
      access_token: jwt
      scope: wildcard
    
    ttl:
      login_request: 1h
      consent_request: 1h
      access_token: 1h
      refresh_token: 24h
      id_token: 1h
      auth_code: 10m
    
    automigration:
      enabled: true
EOF

kubectl apply -f k3s/auth/hydra-deployment.yaml
wait_for_deployment auth hydra

# Deploy authentication handler
echo "🔑 Deploying Authentication Handler..."
kubectl apply -f k3s/auth/auth-handler-deployment.yaml
wait_for_deployment auth auth-handler

# Deploy Oathkeeper
echo "🛡️  Deploying Oathkeeper..."
kubectl apply -f k3s/auth/oathkeeper-deployment.yaml
wait_for_deployment auth oathkeeper

# Deploy backend
echo "🔙 Deploying Backend API..."
kubectl apply -f k3s/backend/backend-deployment.yaml
wait_for_deployment backend backend

# Deploy cloud-build
echo "🏗️  Deploying Cloud Build Service..."
kubectl apply -f k3s/cloud-build/cloud-build-deployment.yaml
wait_for_deployment cloud-build cloud-build

# Step 10: Setup Production Ingress
echo -e "\n${BLUE}Step 10: Setting up Production Ingress${NC}"
echo "----------------------------------------"

cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@$DOMAIN
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: naytife-ingress
  namespace: auth
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
spec:
  tls:
  - hosts:
    - $DOMAIN
    secretName: naytife-tls
  rules:
  - host: $DOMAIN
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: oathkeeper-proxy
            port:
              number: 4455
EOF

echo -e "${GREEN}✅ Production ingress configured with SSL${NC}"

# Step 11: Initialize OAuth2 Clients
echo -e "\n${BLUE}Step 11: Initializing OAuth2 Clients${NC}"
echo "----------------------------------------"

echo "⚙️  Creating OAuth2 clients in Hydra..."

# Wait for Hydra to be fully ready
sleep 60

# Generate secure client secrets
DASHBOARD_CLIENT_SECRET=$(openssl rand -base64 32)
INTROSPECTION_CLIENT_SECRET=$(openssl rand -base64 32)

# Create dashboard client
kubectl exec -n auth deployment/hydra -- \
    hydra create client \
    --endpoint http://127.0.0.1:4445 \
    --id naytife-dashboard \
    --secret $DASHBOARD_CLIENT_SECRET \
    --grant-types authorization_code,refresh_token \
    --response-types code \
    --scope openid,offline,hydra.openid,introspect \
    --callbacks https://$DOMAIN/auth/callback/ory-hydra \
    --name "Naytife Dashboard" || echo "Client may already exist, continuing..."

# Create introspection client
kubectl exec -n auth deployment/hydra -- \
    hydra create client \
    --endpoint http://127.0.0.1:4445 \
    --id naytife-introspection \
    --secret $INTROSPECTION_CLIENT_SECRET \
    --grant-types client_credentials \
    --scope introspect \
    --name "Naytife Introspection" || echo "Client may already exist, continuing..."

echo -e "${GREEN}✅ OAuth2 clients configured${NC}"

# Save client credentials
cat > naytife-oauth-clients.txt <<EOF
=== Naytife OAuth2 Client Credentials ===

Dashboard Client:
  Client ID: naytife-dashboard
  Client Secret: $DASHBOARD_CLIENT_SECRET

Introspection Client:
  Client ID: naytife-introspection  
  Client Secret: $INTROSPECTION_CLIENT_SECRET

Store these credentials securely!
EOF

echo -e "${YELLOW}⚠️  OAuth2 client credentials saved to naytife-oauth-clients.txt${NC}"

# Step 12: Final Verification
echo -e "\n${BLUE}Step 12: Final Verification${NC}"
echo "----------------------------------------"

echo "🧪 Verifying all deployments..."

# Check all deployments
for ns_dep in "database/postgres-postgresql" "database/redis-master" "auth/hydra" "auth/auth-handler" "auth/oathkeeper" "backend/backend" "cloud-build/cloud-build"; do
    ns=$(echo $ns_dep | cut -d'/' -f1)
    dep=$(echo $ns_dep | cut -d'/' -f2)
    
    if kubectl get deployment $dep -n $ns >/dev/null 2>&1; then
        if kubectl wait --for=condition=available deployment/$dep -n $ns --timeout=30s >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $dep ($ns) is ready${NC}"
        else
            echo -e "${YELLOW}⚠️  $dep ($ns) is not ready yet${NC}"
        fi
    else
        echo -e "${RED}❌ $dep ($ns) not found${NC}"
    fi
done

# Check SSL certificate
echo "🔒 Checking SSL certificate..."
sleep 30
if kubectl get certificate naytife-tls -n auth -o jsonpath='{.status.conditions[0].status}' | grep -q True; then
    echo -e "${GREEN}✅ SSL certificate is ready${NC}"
else
    echo -e "${YELLOW}⚠️  SSL certificate is still being provisioned${NC}"
fi

# Step 13: Summary
echo -e "\n${BLUE}Step 13: Production Deployment Complete${NC}"
echo "============================================="

echo -e "${GREEN}🎉 Naytife is now deployed in production!${NC}"
echo ""
echo -e "${PURPLE}🌐 Access Points:${NC}"
echo "  🏠 Main Application:  https://$DOMAIN"
echo "  🔐 OAuth2 Login:      https://$DOMAIN/oauth2/auth"
echo "  🔙 Backend API:       https://$DOMAIN/v1/"
echo "  📚 API Documentation: https://$DOMAIN/v1/docs"
echo "  🏗️  Cloud Build:       https://$DOMAIN/build/"
echo ""
echo -e "${PURPLE}🔑 Important Security Notes:${NC}"
echo "  • OAuth2 client credentials are saved in naytife-oauth-clients.txt"
echo "  • Database password: ${DB_PASSWORD:0:8}..."
echo "  • Redis password: ${REDIS_PASSWORD:0:8}..."
echo "  • Store all credentials securely!"
echo ""
echo -e "${PURPLE}🔧 Management Commands:${NC}"
echo "  kubectl get pods --all-namespaces  # Check all pods"
echo "  kubectl logs -n auth deployment/hydra -f  # View Hydra logs"
echo "  kubectl logs -n backend deployment/backend -f  # View backend logs"
echo "  sudo systemctl status k3s  # Check k3s status"
echo ""
echo -e "${PURPLE}📊 Monitoring:${NC}"
echo "  kubectl top nodes  # Node resource usage"
echo "  kubectl top pods --all-namespaces  # Pod resource usage"
echo ""
echo -e "${GREEN}Production deployment successful! 🚀${NC}"
