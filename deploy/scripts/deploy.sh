#!/bin/bash

# Production-Grade Deploy Script for Kustomize-based deployments
# Usage: ./deploy.sh <environment> [options]

set -euo pipefail

# Global variables
SCRIPT_VERSION="2.0.0"
SCRIPT_NAME="$(basename "$0")"
DEPLOYMENT_ID="deploy-$(date +%Y%m%d_%H%M%S)"
LOG_FILE="/tmp/naytife-deploy-${DEPLOYMENT_ID}.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output with logging
print_info() { echo -e "${BLUE}ℹ️ $1${NC}" | tee -a "$LOG_FILE"; }
print_success() { echo -e "${GREEN}✅ $1${NC}" | tee -a "$LOG_FILE"; }
print_warning() { echo -e "${YELLOW}⚠️ $1${NC}" | tee -a "$LOG_FILE"; }
print_error() { echo -e "${RED}❌ $1${NC}" | tee -a "$LOG_FILE"; }
print_header() { echo -e "${CYAN}📦 $1${NC}" | tee -a "$LOG_FILE"; }
print_debug() { echo -e "${MAGENTA}🔍 $1${NC}" | tee -a "$LOG_FILE"; }

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$LOG_FILE"
}

# Error handling and cleanup
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        print_error "Deployment failed with exit code $exit_code"
        log "DEPLOYMENT FAILED: Exit code $exit_code"
        
        # Show recent events on failure
        if [ -n "${NAMESPACE:-}" ]; then
            print_info "Recent events in namespace $NAMESPACE:"
            kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -n 10 || true
        fi
    fi
    
    # Cleanup temporary files
    rm -f "/tmp/kustomize-output-$ENVIRONMENT.yaml" 2>/dev/null || true
    rm -rf "${TEMP_DIR:-}" 2>/dev/null || true
    
    if [ $exit_code -eq 0 ]; then
        print_success "Deployment logs saved to: $LOG_FILE"
    else
        print_error "Deployment logs saved to: $LOG_FILE"
    fi
}

trap cleanup EXIT

# Enhanced usage function
show_usage() {
    cat << EOF
$SCRIPT_NAME v$SCRIPT_VERSION - Production-Grade Kubernetes Deployment Script

Usage: $SCRIPT_NAME <environment> [options]

Environments:
  local        Deploy to local development environment
  staging      Deploy to staging environment  
  production   Deploy to production environment

Options:
  --dry-run              Show what would be deployed without applying
  --force                Skip confirmation prompts
  --no-secrets          Skip secrets deployment
  --timeout=DURATION    Deployment timeout (default: 600s)
  --wait                Wait for deployment to be ready
  --rollback            Rollback to previous deployment
  --validate            Validate configuration before deployment
  --debug               Enable debug output
  --help, -h            Show this help

Examples:
  $SCRIPT_NAME local                          # Deploy to local
  $SCRIPT_NAME staging --dry-run              # Preview staging deployment
  $SCRIPT_NAME production --force --wait      # Force deploy to production and wait
  $SCRIPT_NAME local --rollback               # Rollback local deployment

EOF
}

# Default values
ENVIRONMENT=""
DRY_RUN=false
FORCE=false
NO_SECRETS=false
TIMEOUT="600s"
WAIT_FOR_READY=false
ROLLBACK=false
VALIDATE_ONLY=false
DEBUG=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        local|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --no-secrets)
            NO_SECRETS=true
            shift
            ;;
        --timeout=*)
            TIMEOUT="${1#*=}"
            shift
            ;;
        --wait)
            WAIT_FOR_READY=true
            shift
            ;;
        --rollback)
            ROLLBACK=true
            shift
            ;;
        --validate)
            VALIDATE_ONLY=true
            shift
            ;;
        --debug)
            DEBUG=true
            set -x
            shift
            ;;
        --help|-h)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check if environment is provided
if [ -z "$ENVIRONMENT" ]; then
    print_error "Environment is required"
    show_usage
    exit 1
fi

log "DEPLOYMENT STARTED: $ENVIRONMENT environment by $(whoami) on $(hostname)"
print_header "Naytife Deployment v$SCRIPT_VERSION - $ENVIRONMENT Environment"
echo "==============================================================="

# Set environment-specific variables
case $ENVIRONMENT in
    local)
        NAMESPACE="naytife"
        CONTEXT_PATTERN="k3d"
        ;;
    staging)
        NAMESPACE="naytife-staging"
        CONTEXT_PATTERN="staging"
        ;;
    production)
        NAMESPACE="naytife-production"
        CONTEXT_PATTERN="production"
        ;;
    *)
        print_error "Invalid environment: $ENVIRONMENT"
        print_info "Available environments: local, staging, production"
        exit 1
        ;;
esac

log "TARGET: $ENVIRONMENT environment, namespace: $NAMESPACE"

# Production safety check
if [ "$ENVIRONMENT" = "production" ] && [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
    print_warning "🚨 PRODUCTION DEPLOYMENT WARNING 🚨"
    print_warning "You are about to deploy to PRODUCTION environment!"
    echo ""
    print_info "Deployment details:"
    print_info "  • Environment: $ENVIRONMENT"
    print_info "  • Namespace: $NAMESPACE"
    print_info "  • Deployment ID: $DEPLOYMENT_ID"
    print_info "  • User: $(whoami)"
    print_info "  • Host: $(hostname)"
    echo ""
    read -p "Are you sure you want to continue? (type 'DEPLOY TO PRODUCTION' to confirm): " confirmation
    
    if [ "$confirmation" != "DEPLOY TO PRODUCTION" ]; then
        print_info "Production deployment cancelled"
        exit 0
    fi
    log "PRODUCTION DEPLOYMENT CONFIRMED by user"
fi

# Set script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
OVERLAY_DIR="$DEPLOY_DIR/overlays/$ENVIRONMENT"

# Check if overlay directory exists
if [ ! -d "$OVERLAY_DIR" ]; then
    print_error "Overlay directory not found: $OVERLAY_DIR"
    exit 1
fi

# Prerequisites check
print_header "Prerequisites Check"

# Check required tools
tools_required=("kubectl" "kustomize")
if [ "$NO_SECRETS" = false ]; then
    tools_required+=("sops")
fi

for tool in "${tools_required[@]}"; do
    if ! command -v "$tool" &> /dev/null; then
        print_error "$tool is not installed"
        exit 1
    fi
    print_success "$tool is available"
done

# Set up SOPS environment if needed
if [ "$NO_SECRETS" = false ]; then
    if [ -z "$SOPS_AGE_KEY_FILE" ]; then
        if [ -f "$HOME/.config/sops/age/keys.txt" ]; then
            export SOPS_AGE_KEY_FILE="$HOME/.config/sops/age/keys.txt"
            print_info "Using SOPS age key file: $SOPS_AGE_KEY_FILE"
        else
            print_warning "SOPS age key file not found at $HOME/.config/sops/age/keys.txt"
            print_info "Make sure to set SOPS_AGE_KEY_FILE environment variable if using a custom location"
        fi
    else
        print_info "Using SOPS age key file: $SOPS_AGE_KEY_FILE"
    fi
fi

# Check cluster connectivity
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

current_context=$(kubectl config current-context)
print_success "Connected to cluster: $current_context"

# Verify cluster context for environment
if [ "$ENVIRONMENT" != "local" ]; then
    if ! echo "$current_context" | grep -q "$CONTEXT_PATTERN"; then
        print_warning "Current context '$current_context' may not match $ENVIRONMENT environment"
        if [ "$FORCE" = false ]; then
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                print_info "Deployment cancelled"
                exit 0
            fi
        fi
    fi
fi

log "CLUSTER: $current_context"

# Ensure namespace exists
print_info "Checking namespace..."
if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
    print_warning "Namespace '$NAMESPACE' does not exist"
    if [ "$DRY_RUN" = false ]; then
        print_info "Creating namespace '$NAMESPACE'"
        kubectl create namespace "$NAMESPACE"
        print_success "Namespace created"
    else
        print_info "Would create namespace '$NAMESPACE'"
    fi
else
    print_success "Namespace '$NAMESPACE' exists"
fi

# Ensure auth namespace exists for environments that need it
AUTH_NAMESPACE="${NAMESPACE}-auth"
print_info "Checking auth namespace..."
if ! kubectl get namespace "$AUTH_NAMESPACE" >/dev/null 2>&1; then
    print_warning "Auth namespace '$AUTH_NAMESPACE' does not exist"
    if [ "$DRY_RUN" = false ]; then
        print_info "Creating auth namespace '$AUTH_NAMESPACE'"
        kubectl create namespace "$AUTH_NAMESPACE"
        print_success "Auth namespace created"
    else
        print_info "Would create auth namespace '$AUTH_NAMESPACE'"
    fi
else
    print_success "Auth namespace '$AUTH_NAMESPACE' exists"
fi

# Handle rollback
if [ "$ROLLBACK" = true ]; then
    print_header "Rolling Back Deployment"
    
    # Check if there are previous revisions
    deployments=$(kubectl get deployments -n "$NAMESPACE" --no-headers -o name 2>/dev/null || true)
    if [ -z "$deployments" ]; then
        print_error "No deployments found in namespace '$NAMESPACE'"
        exit 1
    fi
    
    print_info "Rolling back deployments in namespace '$NAMESPACE'..."
    echo "$deployments" | while read deployment; do
        deployment_name=$(basename "$deployment")
        print_info "Rolling back $deployment_name"
        kubectl rollout undo "$deployment" -n "$NAMESPACE"
    done
    
    print_success "Rollback initiated for all deployments"
    exit 0
fi

# Configuration validation
print_header "Configuration Validation"

# Change to overlay directory
cd "$OVERLAY_DIR"

print_info "Building Kustomize configuration..."
log "KUSTOMIZE BUILD: $OVERLAY_DIR"

# Build the configuration with validation
if ! kustomize build . > "/tmp/kustomize-output-$ENVIRONMENT.yaml"; then
    print_error "Failed to build Kustomize configuration"
    exit 1
fi

print_success "Configuration built successfully"

# Validate configuration
if command -v kubeval &> /dev/null; then
    print_info "Validating Kubernetes manifests with kubeval..."
    if kubeval "/tmp/kustomize-output-$ENVIRONMENT.yaml"; then
        print_success "Manifest validation passed"
    else
        print_warning "Manifest validation warnings found"
    fi
fi

# Show configuration preview if debug mode
if [ "$DEBUG" = true ]; then
    print_debug "Configuration preview:"
    echo "--- Configuration START ---"
    cat "/tmp/kustomize-output-$ENVIRONMENT.yaml"
    echo "--- Configuration END ---"
fi

# Exit if validate-only mode
if [ "$VALIDATE_ONLY" = true ]; then
    print_success "Configuration validation completed"
    exit 0
fi

# Secrets management
if [ "$NO_SECRETS" = false ]; then
    print_header "Secrets Management"
    
    if command -v sops &> /dev/null; then
        print_info "Processing secrets with SOPS..."
        
        # Create temporary directory for decrypted secrets
        TEMP_DIR=$(mktemp -d)
        
        secret_files_found=false
        
        # Check if secrets directory exists
        if [ -d "$DEPLOY_DIR/secrets/$ENVIRONMENT" ]; then
            # Decrypt secrets and apply them
            for secret_file in "$DEPLOY_DIR/secrets/$ENVIRONMENT"/*.yaml; do
                if [ -f "$secret_file" ]; then
                    secret_files_found=true
                    filename=$(basename "$secret_file")
                    print_info "Processing secret: $filename"
                    log "DECRYPT SECRET: $filename"
                    
                    if sops -d "$secret_file" > "$TEMP_DIR/$filename"; then
                        print_success "Decrypted $filename"
                    else
                        print_error "Failed to decrypt $filename"
                        exit 1
                    fi
                fi
            done
        fi
        
        if [ "$secret_files_found" = true ]; then
            # Apply decrypted secrets
            if [ "$DRY_RUN" = true ]; then
                print_info "Secrets that would be applied (dry-run mode):"
                for secret_file in "$TEMP_DIR"/*.yaml; do
                    if [ -f "$secret_file" ]; then
                        echo "--- Secret: $(basename "$secret_file") ---"
                        # Show metadata only for security
                        kubectl apply --dry-run=client -f "$secret_file" -o yaml | head -n 10
                    fi
                done
            else
                print_info "Applying secrets..."
                for secret_file in "$TEMP_DIR"/*.yaml; do
                    if [ -f "$secret_file" ]; then
                        print_info "Applying $(basename "$secret_file")"
                        kubectl apply -f "$secret_file"
                    fi
                done
                print_success "Secrets applied successfully"
            fi
        else
            print_warning "No secret files found in $DEPLOY_DIR/secrets/$ENVIRONMENT"
        fi
    else
        print_warning "SOPS not available. Secrets will not be processed."
        if [ "$ENVIRONMENT" != "local" ]; then
            print_error "SOPS is required for $ENVIRONMENT environment"
            exit 1
        fi
    fi
else
    print_info "Skipping secrets deployment (--no-secrets flag)"
fi

# Main deployment
print_header "Main Deployment"

if [ "$DRY_RUN" = true ]; then
    print_info "Configuration that would be applied (DRY RUN):"
    echo "================================================"
    cat "/tmp/kustomize-output-$ENVIRONMENT.yaml"
    echo "================================================"
    print_success "Dry run completed successfully"
else
    print_info "Applying configuration to Kubernetes cluster..."
    log "APPLY CONFIG: Starting main deployment"
    
    # Apply with server-side validation
    if kubectl apply -f "/tmp/kustomize-output-$ENVIRONMENT.yaml" --validate=true; then
        print_success "Configuration applied successfully!"
        log "APPLY CONFIG: Success"
        
        # Wait for deployment readiness if requested
        if [ "$WAIT_FOR_READY" = true ]; then
            print_header "Waiting for Deployment Readiness"
            
            print_info "Waiting for deployments to be ready (timeout: $TIMEOUT)..."
            
            # Get all deployments that were just applied
            applied_deployments=$(kubectl get deployments -n "$NAMESPACE" --no-headers -o name 2>/dev/null || true)
            
            if [ -n "$applied_deployments" ]; then
                # Wait for each deployment
                echo "$applied_deployments" | while read deployment; do
                    deployment_name=$(basename "$deployment")
                    print_info "Waiting for $deployment_name..."
                    
                    if kubectl rollout status "$deployment" -n "$NAMESPACE" --timeout="$TIMEOUT"; then
                        print_success "$deployment_name is ready"
                    else
                        print_error "$deployment_name failed to become ready within $TIMEOUT"
                        return 1
                    fi
                done
                
                if [ $? -eq 0 ]; then
                    print_success "All deployments are ready!"
                    log "READINESS: All deployments ready"
                else
                    print_error "Some deployments failed to become ready"
                    log "READINESS: Failed"
                    exit 1
                fi
            else
                print_warning "No deployments found to wait for"
            fi
        fi
        
        # Post-deployment status check
        print_header "Post-Deployment Status"
        
        print_info "Deployment status:"
        kubectl get deployments -n "$NAMESPACE" -o wide || true
        
        print_info "Pod status:"
        kubectl get pods -n "$NAMESPACE" -o wide || true
        
        print_info "Service status:"
        kubectl get services -n "$NAMESPACE" -o wide || true
        
        # Check for any failed pods
        failed_pods=$(kubectl get pods -n "$NAMESPACE" --field-selector=status.phase=Failed --no-headers 2>/dev/null | wc -l)
        if [ "$failed_pods" -gt 0 ]; then
            print_warning "$failed_pods pod(s) in Failed state"
            kubectl get pods -n "$NAMESPACE" --field-selector=status.phase=Failed
        fi
        
        # Show recent events
        print_info "Recent events:"
        kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -n 5 || true
        
    else
        print_error "Deployment failed!"
        log "APPLY CONFIG: Failed"
        exit 1
    fi
fi

# Deployment completion
print_header "Deployment Summary"

log "DEPLOYMENT COMPLETED: $ENVIRONMENT environment"

if [ "$DRY_RUN" = true ]; then
    print_success "Dry run completed successfully!"
    print_info "No changes were applied to the cluster"
else
    print_success "Deployment completed successfully! 🎉"
    print_info "Deployment ID: $DEPLOYMENT_ID"
    print_info "Environment: $ENVIRONMENT"
    print_info "Namespace: $NAMESPACE"
    print_info "Deployed by: $(whoami) on $(hostname)"
    print_info "Timestamp: $(date)"
fi

# Environment-specific post-deployment instructions
case $ENVIRONMENT in
    local)
        print_header "Local Environment Access"
        print_info "Services are accessible via NodePort or port-forward"
        print_info "Use './status.sh local' to check deployment status"
        print_info "Use './logs.sh [service] local' to view service logs"
        print_info "Use './test-integration.sh local' to run tests"
        
        # Show NodePort services
        nodeport_services=$(kubectl get services -n "$NAMESPACE" -o jsonpath='{.items[?(@.spec.type=="NodePort")].metadata.name}' 2>/dev/null || true)
        if [ -n "$nodeport_services" ]; then
            print_info "NodePort services available:"
            for service in $nodeport_services; do
                port=$(kubectl get service "$service" -n "$NAMESPACE" -o jsonpath='{.spec.ports[0].nodePort}')
                print_info "  • $service: http://localhost:$port"
            done
        fi
        ;;
    staging)
        print_header "Staging Environment Access"
        print_info "Services should be accessible via staging URLs"
        print_info "Monitor deployment: './status.sh staging'"
        print_info "Run integration tests: './test-integration.sh staging'"
        
        # Show ingress URLs if available
        ingress_hosts=$(kubectl get ingress -n "$NAMESPACE" -o jsonpath='{.items[*].spec.rules[*].host}' 2>/dev/null || true)
        if [ -n "$ingress_hosts" ]; then
            print_info "Ingress endpoints:"
            for host in $ingress_hosts; do
                print_info "  • https://$host"
            done
        fi
        ;;
    production)
        print_header "Production Environment Access"
        print_warning "🚨 PRODUCTION DEPLOYMENT COMPLETED 🚨"
        print_info "Monitor the deployment closely for any issues"
        print_info "Health checks: './status.sh production'"
        print_info "Performance tests: './benchmark.sh production'"
        
        # Show production URLs
        print_info "Production endpoints:"
        print_info "  • API: https://api.naytife.com"
        print_info "  • Auth: https://auth.naytife.com"
        
        # Production monitoring reminders
        print_warning "IMPORTANT: Monitor the following:"
        print_warning "  • Application logs and metrics"
        print_warning "  • Response times and error rates"
        print_warning "  • Resource utilization"
        print_warning "  • User impact and feedback"
        ;;
esac

print_success "All operations completed successfully!"
echo ""
print_info "Next steps:"
print_info "  • Review deployment status: './status.sh $ENVIRONMENT'"
print_info "  • Run integration tests: './test-integration.sh $ENVIRONMENT'"
print_info "  • Monitor logs: './logs.sh [service] $ENVIRONMENT'"
if [ "$ENVIRONMENT" = "local" ]; then
    print_info "  • Run benchmarks: './benchmark.sh $ENVIRONMENT'"
fi

exit 0
