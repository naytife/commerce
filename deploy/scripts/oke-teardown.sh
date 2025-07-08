#!/bin/bash

# Teardown Naytife Platform from OKE Always Free Cluster
# This script safely removes the application while preserving data if needed

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

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
OVERLAY_DIR="$DEPLOY_DIR/overlays/$ENVIRONMENT"

print_header "Tearing Down Naytife Platform from OKE Always Free Cluster"

# Set kubeconfig
export KUBECONFIG="$OKE_KUBECONFIG"

# Check if connected to cluster
if ! kubectl cluster-info &>/dev/null; then
    print_error "Cannot connect to OKE cluster. Please check your kubeconfig."
    exit 1
fi

print_success "Connected to OKE cluster"

# Check if namespace exists
if ! kubectl get namespace "$NAMESPACE" &>/dev/null; then
    print_warning "Namespace $NAMESPACE does not exist. Nothing to tear down."
    exit 0
fi

# Ask for confirmation
print_warning "This will delete all resources in namespace: $NAMESPACE"
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Teardown cancelled."
    exit 0
fi

# Option to preserve data
print_info "Do you want to preserve persistent volumes (database data)?"
read -p "Preserve data? (Y/n): " -n 1 -r
echo
PRESERVE_DATA=true
if [[ $REPLY =~ ^[Nn]$ ]]; then
    PRESERVE_DATA=false
fi

# Show current resources
print_header "Current Resources in $NAMESPACE"
kubectl get all,pvc,ingress -n "$NAMESPACE"

# Delete application resources
print_header "Deleting Application Resources"

if [ -d "$OVERLAY_DIR" ]; then
    print_info "Using Kustomize to delete resources..."
    kustomize build "$OVERLAY_DIR" | kubectl delete -f - --ignore-not-found=true
else
    print_info "Deleting resources directly..."
    kubectl delete all --all -n "$NAMESPACE" --ignore-not-found=true
    kubectl delete ingress --all -n "$NAMESPACE" --ignore-not-found=true
fi

print_success "Application resources deleted"

# Handle persistent volumes
if [ "$PRESERVE_DATA" = true ]; then
    print_info "Preserving persistent volumes and claims"
    print_warning "PVCs will remain in the namespace. To delete them later:"
    echo "kubectl delete pvc --all -n $NAMESPACE"
else
    print_info "Deleting persistent volume claims..."
    kubectl delete pvc --all -n "$NAMESPACE" --ignore-not-found=true
    print_success "Persistent volume claims deleted"
fi

# Option to delete namespace
print_info "Do you want to delete the entire namespace?"
print_warning "This will remove everything including any remaining PVCs"
read -p "Delete namespace $NAMESPACE? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "Deleting namespace..."
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
    print_success "Namespace deleted"
else
    print_info "Namespace preserved"
fi

# Clean up storage class (optional)
print_info "Do you want to remove the OCI storage class?"
print_warning "This might affect other applications using the same storage class"
read -p "Delete OCI storage class? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    kubectl delete storageclass oci-bv --ignore-not-found=true
    print_success "Storage class deleted"
fi

print_header "Teardown Complete"

if [ "$PRESERVE_DATA" = true ]; then
    print_info "Data preservation summary:"
    echo "  - Application pods: DELETED"
    echo "  - Services and ingress: DELETED"
    echo "  - Persistent volumes: PRESERVED"
    echo "  - Database data: PRESERVED"
    print_info "To restore the application with existing data:"
    echo "  ./oke-deploy.sh"
else
    print_success "Complete teardown finished"
    print_info "All application resources and data have been removed"
fi

print_info "OKE cluster and nodes remain running (Always Free resources)"
print_info "To stop paying for any additional resources, check OCI Console"
