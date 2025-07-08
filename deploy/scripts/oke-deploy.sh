#!/bin/bash

# Deploy Naytife Platform to OKE Always Free Cluster
# This script deploys the application using the existing deployment patterns

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}📦 $1${NC}"; }

# Configuration
ENVIRONMENT="oke-free"
NAMESPACE="naytife-oke-free"
OKE_KUBECONFIG="$HOME/.kube/config-oke"
TIMEOUT="600s"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
OVERLAY_DIR="$DEPLOY_DIR/overlays/$ENVIRONMENT"

print_header "Deploying Naytife Platform to OKE Always Free Cluster"

# Source existing deployment utilities if available
if [ -f "$SCRIPT_DIR/deploy.sh" ]; then
    print_info "Loading existing deployment utilities"
    # Note: We're not sourcing the full script, just using patterns
fi

# Check prerequisites
print_info "Checking prerequisites..."

# Check if kubeconfig exists
if [ ! -f "$OKE_KUBECONFIG" ]; then
    print_error "OKE kubeconfig not found at $OKE_KUBECONFIG"
    print_info "Please run: oci ce cluster create-kubeconfig --cluster-id <cluster-id> --file $OKE_KUBECONFIG"
    exit 1
fi

# Set kubeconfig
export KUBECONFIG="$OKE_KUBECONFIG"

# Test cluster connection
if ! kubectl cluster-info &>/dev/null; then
    print_error "Cannot connect to OKE cluster. Please check your kubeconfig and cluster status."
    exit 1
fi

print_success "Connected to OKE cluster"

# Show cluster info
print_info "Cluster information:"
kubectl get nodes -o wide

# Check if overlay directory exists
if [ ! -d "$OVERLAY_DIR" ]; then
    print_error "Overlay directory not found: $OVERLAY_DIR"
    exit 1
fi

print_success "Found overlay directory: $OVERLAY_DIR"

# Check if kustomization builds correctly
print_info "Validating Kustomize configuration..."
if ! kustomize build "$OVERLAY_DIR" >/dev/null; then
    print_error "Kustomize build failed. Please check your overlay configuration."
    exit 1
fi

print_success "Kustomize configuration is valid"

# Create namespace if it doesn't exist
print_info "Ensuring namespace exists..."
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Apply storage class first
print_info "Applying storage class..."
kubectl apply -f "$OVERLAY_DIR/storage-class.yaml"

# Deploy application using kustomize
print_header "Deploying application stack..."
print_info "Building and applying Kustomize configuration..."

kustomize build "$OVERLAY_DIR" | kubectl apply -f -

print_success "Application manifests applied"

# Wait for deployments to be ready
print_header "Waiting for deployments to be ready..."

DEPLOYMENTS=("postgres" "redis" "hydra" "oathkeeper" "backend" "auth-handler" "template-registry" "store-deployer")

for deployment in "${DEPLOYMENTS[@]}"; do
    print_info "Waiting for deployment: $deployment"
    if kubectl rollout status deployment/"$deployment" -n "$NAMESPACE" --timeout="$TIMEOUT"; then
        print_success "✅ $deployment is ready"
    else
        print_warning "⚠️ $deployment failed to become ready within timeout"
    fi
done

# Show deployment status
print_header "Deployment Status"
kubectl get deployments -n "$NAMESPACE" -o wide

# Show pod status
print_header "Pod Status"
kubectl get pods -n "$NAMESPACE" -o wide

# Show services
print_header "Services"
kubectl get services -n "$NAMESPACE" -o wide

# Show ingress
print_header "Ingress Configuration"
kubectl get ingress -n "$NAMESPACE" -o wide

# Get load balancer IP
LB_IP=$(kubectl get ingress naytife-ingress -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
if [ -n "$LB_IP" ]; then
    print_success "Load Balancer IP: $LB_IP"
    print_info "Add these entries to your /etc/hosts file for local testing:"
    echo "$LB_IP api.naytife-oke.dev"
    echo "$LB_IP auth.naytife-oke.dev"
    echo "$LB_IP oauth.naytife-oke.dev"
    echo "$LB_IP gateway.naytife-oke.dev"
else
    print_info "Load Balancer IP not yet assigned. Check again in a few minutes:"
    echo "kubectl get ingress naytife-ingress -n $NAMESPACE"
fi

# Show persistent volumes
print_header "Persistent Volumes"
kubectl get pv,pvc -n "$NAMESPACE"

# Show resource usage
if kubectl top nodes &>/dev/null; then
    print_header "Resource Usage"
    print_info "Node resource usage:"
    kubectl top nodes
    print_info "Pod resource usage:"
    kubectl top pods -n "$NAMESPACE"
fi

print_header "Deployment Complete!"
print_success "Naytife platform deployed to OKE Always Free cluster"

print_info "Useful commands:"
echo "  View logs: kubectl logs -f deployment/<service-name> -n $NAMESPACE"
echo "  Shell into pod: kubectl exec -it deployment/<service-name> -n $NAMESPACE -- /bin/sh"
echo "  Port forward: kubectl port-forward service/<service-name> <local-port>:<service-port> -n $NAMESPACE"
echo "  Scale deployment: kubectl scale deployment <service-name> --replicas=<count> -n $NAMESPACE"

print_info "To update the deployment:"
echo "  1. Make changes to overlay files"
echo "  2. Run: kustomize build $OVERLAY_DIR | kubectl apply -f -"

print_info "To tear down:"
echo "  ./oke-teardown.sh"
