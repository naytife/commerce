#!/bin/bash

# OKE Always Free Cluster Setup Script
# This script helps set up the prerequisites for deploying to OKE Always Free

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() { echo -e "${BLUE}ℹ️ $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_header() { echo -e "${CYAN}📦 $1${NC}"; }

# Configuration
OKE_KUBECONFIG="$HOME/.kube/config-oke"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

print_header "OKE Always Free Cluster Setup"

# Check if OCI CLI is installed
if ! command -v oci &> /dev/null; then
    print_error "OCI CLI is not installed. Please install it first:"
    echo "curl -L https://raw.githubusercontent.com/oracle/oci-cli/master/scripts/install/install.sh | bash"
    exit 1
fi

print_success "OCI CLI found"

# Check if OCI CLI is configured
print_info "Testing OCI CLI authentication..."
if oci os ns get >/dev/null 2>&1; then
    print_success "OCI CLI is configured and authenticated"
elif [ ! -f "$HOME/.oci/config" ]; then
    print_error "OCI CLI is not configured. Please run: oci setup config"
    exit 1
else
    print_warning "OCI CLI is configured but authentication is failing"
    print_info "This could be due to:"
    echo "  1. API key not yet active in OCI Console (can take a few minutes)"
    echo "  2. Fingerprint mismatch between local and OCI Console"
    echo "  3. API key not properly uploaded to OCI Console"
    echo "  4. Missing IAM permissions - your user needs to be in a group with policies"
    echo ""
    print_info "Current fingerprint: $(grep fingerprint ~/.oci/config | cut -d='=' -f2)"
    print_info "Please verify this matches exactly in OCI Console"
    echo ""
    print_error "REQUIRED IAM SETUP:"
    echo "1. Create group 'OKEDevelopers' in OCI Console → Identity & Security → Groups"
    echo "2. Add your user to this group"
    echo "3. Create policy with these statements:"
    echo "   Allow group OKEDevelopers to manage cluster-family in tenancy"
    echo "   Allow group OKEDevelopers to manage compute-instances in tenancy"
    echo "   Allow group OKEDevelopers to manage repos in tenancy"
    echo "   Allow group OKEDevelopers to read objectstorage-namespaces in tenancy"
    echo "   Allow group OKEDevelopers to manage object-family in tenancy"
    echo "   Allow group OKEDevelopers to manage vnics in tenancy"
    echo "   Allow group OKEDevelopers to manage subnets in tenancy"
    echo "   Allow group OKEDevelopers to manage vcns in tenancy"
    echo "   (Or simply: Allow group OKEDevelopers to manage all-resources in tenancy)"
    echo ""
    print_info "To test authentication manually: oci os ns get"
    print_info "We'll continue with setup, but authentication must work before deployment"
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install it first:"
    echo "https://kubernetes.io/docs/tasks/tools/"
    exit 1
fi

print_success "kubectl found"

# Check if kustomize is installed
if ! command -v kustomize &> /dev/null; then
    print_error "kustomize is not installed. Please install it first:"
    echo "https://kubectl.docs.kubernetes.io/installation/kustomize/"
    exit 1
fi

print_success "kustomize found"

# Check if docker is installed and running
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker."
    exit 1
fi

print_success "Docker is running"

# Check if docker buildx is available for ARM64 builds
if ! docker buildx version &> /dev/null; then
    print_error "Docker buildx is not available. Please update Docker to a recent version."
    exit 1
fi

print_success "Docker buildx available for multi-architecture builds"

# Check OKE kubeconfig
if [ -f "$OKE_KUBECONFIG" ]; then
    print_info "Found existing OKE kubeconfig at $OKE_KUBECONFIG"
    
    # Test connection
    if KUBECONFIG="$OKE_KUBECONFIG" kubectl cluster-info &>/dev/null; then
        print_success "Successfully connected to OKE cluster"
        
        # Show cluster info
        print_info "Cluster information:"
        KUBECONFIG="$OKE_KUBECONFIG" kubectl get nodes -o wide
        
        # Check for ARM nodes
        ARM_NODES=$(KUBECONFIG="$OKE_KUBECONFIG" kubectl get nodes -o jsonpath='{.items[*].status.nodeInfo.architecture}' 2>/dev/null | grep -c arm64 2>/dev/null || echo "0")
        if [ "$ARM_NODES" -gt 0 ]; then
            print_success "Found $ARM_NODES ARM64 node(s) in cluster"
        else
            print_warning "No ARM64 nodes found yet. Node pools may still be starting up."
        fi
    else
        print_warning "Cannot connect to OKE cluster. Please check your kubeconfig."
        print_info "To create/update kubeconfig, run:"
        echo "oci ce cluster create-kubeconfig --cluster-id <cluster-id> --file $OKE_KUBECONFIG --region <region> --token-version 2.0.0"
    fi
else
    print_warning "No OKE kubeconfig found at $OKE_KUBECONFIG"
    print_info "To create kubeconfig after setting up OKE cluster, run:"
    echo "oci ce cluster create-kubeconfig --cluster-id <cluster-id> --file $OKE_KUBECONFIG --region <region> --token-version 2.0.0"
fi

print_header "OCI Container Registry (OCIR) Setup"

# Get user information for OCIR
TENANCY_NAMESPACE=$(oci os ns get --query 'data' --raw-output 2>/dev/null || echo "")
if [ -n "$TENANCY_NAMESPACE" ]; then
    print_success "Tenancy namespace: $TENANCY_NAMESPACE"
    
    # Get current region
    REGION=$(oci iam region-subscription list --query 'data[?"is-home-region"]."region-name"' --raw-output 2>/dev/null | grep -o 'us-[a-z0-9-]*\|eu-[a-z0-9-]*\|ap-[a-z0-9-]*\|ca-[a-z0-9-]*\|uk-[a-z0-9-]*' | head -1 || grep region ~/.oci/config | cut -d'=' -f2 | tr -d ' ' || echo "")
    if [ -n "$REGION" ]; then
        print_success "Current region: $REGION"
        
        OCIR_ENDPOINT="$REGION.ocir.io"
        print_info "OCIR endpoint: $OCIR_ENDPOINT"
        print_info "Your image names should be: $OCIR_ENDPOINT/$TENANCY_NAMESPACE/naytife/<service>:arm64-latest"
    else
        print_warning "Could not determine current region"
    fi
else
    print_warning "Could not determine tenancy namespace"
fi

print_header "Next Steps"

if [ -f "$OKE_KUBECONFIG" ] && KUBECONFIG="$OKE_KUBECONFIG" kubectl cluster-info &>/dev/null; then
    print_info "Your environment is ready! Next steps:"
    echo "1. Update deploy/overlays/oke-free/kustomization.yaml with your OCIR details"
    echo "2. Build and push ARM64 images: ./oke-build-push.sh"
    echo "3. Deploy to OKE: ./oke-deploy.sh"
else
    print_info "Please complete OKE cluster setup:"
    echo "1. Create OKE cluster in OCI Console (use Always Free ARM instances)"
    echo "2. Create kubeconfig: oci ce cluster create-kubeconfig ..."
    echo "3. Run this script again to verify setup"
fi

print_info "For detailed instructions, see:"
echo "- OCI_DEPLOYMENT_PLAN.md"
echo "- OCI_ALWAYS_FREE_SERVICES.md"
