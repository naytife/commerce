#!/bin/bash

# Enhanced Image Build Script with Change Detection
# Builds all service images only when changes are detected in their respective codebases

set -euo pipefail

# Script metadata
SCRIPT_VERSION="1.0.0"
SCRIPT_NAME="$(basename "$0")"
BUILD_ID="build-$(date +%Y%m%d_%H%M%S)"
LOG_FILE="/tmp/naytife-build-${BUILD_ID}.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Logging functions
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
        print_error "Build failed with exit code $exit_code"
        log "BUILD FAILED: Exit code $exit_code"
    fi
    
    if [ $exit_code -eq 0 ]; then
        print_success "Build logs saved to: $LOG_FILE"
    else
        print_error "Build logs saved to: $LOG_FILE"
    fi
}

trap cleanup EXIT

# Usage function
show_usage() {
    cat << EOF
$SCRIPT_NAME v$SCRIPT_VERSION - Enhanced Docker Image Build Script

Usage: $SCRIPT_NAME [options]

Options:
  --force              Force rebuild all images regardless of changes
  --service=NAME       Build only specific service (backend,auth-handler,store-deployer,template-registry,migrations)
  --no-cache          Build without Docker cache
  --parallel          Build images in parallel (experimental)
  --registry=URL      Push to specific registry (default: local only)
  --tag=TAG           Use specific tag (default: latest)
  --check-only        Only check for changes, don't build
  --no-k3s-import     Skip importing images to k3s (default: auto-import for local)
  --verbose           Enable verbose output
  --help, -h          Show this help

Examples:
  $SCRIPT_NAME                                    # Build changed services only
  $SCRIPT_NAME --force                           # Force rebuild all services
  $SCRIPT_NAME --service=backend                 # Build only backend service
  $SCRIPT_NAME --registry=registry.example.com  # Build and push to registry
  $SCRIPT_NAME --check-only                      # Check for changes without building

EOF
}

# Default values
FORCE_BUILD=false
SPECIFIC_SERVICE=""
NO_CACHE=false
PARALLEL_BUILD=false
REGISTRY=""
TAG="latest"
CHECK_ONLY=false
NO_K3S_IMPORT=false
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --force)
            FORCE_BUILD=true
            shift
            ;;
        --service=*)
            SPECIFIC_SERVICE="${1#*=}"
            shift
            ;;
        --no-cache)
            NO_CACHE=true
            shift
            ;;
        --parallel)
            PARALLEL_BUILD=true
            shift
            ;;
        --registry=*)
            REGISTRY="${1#*=}"
            shift
            ;;
        --tag=*)
            TAG="${1#*=}"
            shift
            ;;
        --check-only)
            CHECK_ONLY=true
            shift
            ;;
        --no-k3s-import)
            NO_K3S_IMPORT=true
            shift
            ;;
        --verbose)
            VERBOSE=true
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

# Set script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BUILD_CACHE_DIR="$PROJECT_ROOT/.build-cache"

# Create build cache directory
mkdir -p "$BUILD_CACHE_DIR"

print_header "Naytife Platform Image Builder v$SCRIPT_VERSION"
print_info "Build ID: $BUILD_ID"
print_info "Project root: $PROJECT_ROOT"

# Check prerequisites
print_header "Prerequisites Check"

if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed or not in PATH"
    exit 1
fi
print_success "Docker is available"

if ! command -v git &> /dev/null; then
    print_error "Git is not installed or not in PATH"
    exit 1
fi
print_success "Git is available"

# Verify Docker daemon is running
if ! docker info &> /dev/null; then
    print_error "Docker daemon is not running"
    exit 1
fi
print_success "Docker daemon is running"

# Service definitions
declare -A SERVICES=(
    ["backend"]="backend"
    ["backend-migrations"]="backend"
    ["auth-handler"]="auth/authentication-handler"
    ["store-deployer"]="services/store-deployer"
    ["template-registry"]="services/template-registry"
)

declare -A DOCKERFILES=(
    ["backend"]="Dockerfile"
    ["backend-migrations"]="migrations.Dockerfile"
    ["auth-handler"]="Dockerfile"
    ["store-deployer"]="Dockerfile"
    ["template-registry"]="Dockerfile"
)

declare -A IMAGE_NAMES=(
    ["backend"]="naytife/backend"
    ["backend-migrations"]="naytife/backend-migrations"
    ["auth-handler"]="naytife/auth-handler"
    ["store-deployer"]="naytife/store-deployer"
    ["template-registry"]="naytife/template-registry"
)

# Function to get the last commit hash that affected a directory
get_last_commit() {
    local dir="$1"
    git -C "$PROJECT_ROOT" log -1 --format="%H" -- "$dir" 2>/dev/null || echo "unknown"
}

# Function to check if service needs rebuild
needs_rebuild() {
    local service="$1"
    local service_dir="$PROJECT_ROOT/${SERVICES[$service]}"
    local cache_file="$BUILD_CACHE_DIR/${service}.lastbuild"
    
    if [ "$FORCE_BUILD" = true ]; then
        return 0  # Force rebuild
    fi
    
    if [ ! -f "$cache_file" ]; then
        return 0  # First build
    fi
    
    local last_build_commit
    last_build_commit=$(cat "$cache_file" 2>/dev/null || echo "")
    local current_commit
    current_commit=$(get_last_commit "${SERVICES[$service]}")
    
    if [ "$last_build_commit" != "$current_commit" ]; then
        return 0  # Changes detected
    fi
    
    # Check if Docker image exists
    local image_tag="${IMAGE_NAMES[$service]}:$TAG"
    if [ -n "$REGISTRY" ]; then
        image_tag="$REGISTRY/${IMAGE_NAMES[$service]}:$TAG"
    fi
    
    if ! docker image inspect "$image_tag" &> /dev/null; then
        return 0  # Image doesn't exist
    fi
    
    return 1  # No rebuild needed
}

# Function to import Docker image to k3s
import_to_k3s() {
    local image_tag="$1"
    
    # Skip if explicitly disabled
    if [ "$NO_K3S_IMPORT" = true ]; then
        return 0
    fi
    
    # Check if k3s is available
    if ! command -v k3s &> /dev/null; then
        return 0  # Skip if k3s not available
    fi
    
    # Check if we're in a local development environment
    if [ -n "$REGISTRY" ]; then
        return 0  # Skip if pushing to registry
    fi
    
    print_info "Importing $image_tag to k3s..."
    log "K3S IMPORT START: $image_tag"
    
    # Export image from Docker and import to k3s (need sudo for containerd socket access)
    if docker save "$image_tag" | sudo k3s ctr images import -; then
        print_success "Imported $image_tag to k3s"
        log "K3S IMPORT SUCCESS: $image_tag"
    else
        print_warning "Failed to import $image_tag to k3s (continuing anyway)"
        log "K3S IMPORT FAILED: $image_tag"
        # Don't fail the build for k3s import failure
    fi
}

# Function to build a single service
build_service() {
    local service="$1"
    local service_dir="$PROJECT_ROOT/${SERVICES[$service]}"
    local dockerfile="${DOCKERFILES[$service]}"
    local image_name="${IMAGE_NAMES[$service]}"
    local image_tag="$image_name:$TAG"
    
    if [ -n "$REGISTRY" ]; then
        image_tag="$REGISTRY/$image_name:$TAG"
    fi
    
    print_header "Building $service"
    
    if [ ! -d "$service_dir" ]; then
        print_error "Service directory not found: $service_dir"
        return 1
    fi
    
    if [ ! -f "$service_dir/$dockerfile" ]; then
        print_error "Dockerfile not found: $service_dir/$dockerfile"
        return 1
    fi
    
    # Check if rebuild is needed
    if ! needs_rebuild "$service"; then
        print_info "No changes detected for $service, skipping build"
        return 0
    fi
    
    if [ "$CHECK_ONLY" = true ]; then
        print_warning "$service needs rebuild (changes detected)"
        return 0
    fi
    
    print_info "Building image: $image_tag"
    log "BUILD START: $service -> $image_tag"
    
    # Prepare build arguments
    local build_args=()
    build_args+=("--file" "$dockerfile")
    build_args+=("--tag" "$image_tag")
    
    if [ "$NO_CACHE" = true ]; then
        build_args+=("--no-cache")
    fi
    
    # Add build context
    build_args+=(".")
    
    # Execute build
    local start_time
    start_time=$(date +%s)
    
    if (cd "$service_dir" && docker build "${build_args[@]}"); then
        local end_time
        end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        print_success "Built $service in ${duration}s"
        log "BUILD SUCCESS: $service (${duration}s)"
        
        # Import image to k3s if available (for local development)
        import_to_k3s "$image_tag"
        
        # Update build cache
        get_last_commit "${SERVICES[$service]}" > "$BUILD_CACHE_DIR/${service}.lastbuild"
        
        # Push to registry if specified
        if [ -n "$REGISTRY" ]; then
            print_info "Pushing $image_tag to registry..."
            if docker push "$image_tag"; then
                print_success "Pushed $image_tag to registry"
                log "PUSH SUCCESS: $image_tag"
            else
                print_error "Failed to push $image_tag"
                log "PUSH FAILED: $image_tag"
                return 1
            fi
        fi
        
        return 0
    else
        print_error "Failed to build $service"
        log "BUILD FAILED: $service"
        return 1
    fi
}

# Function to build services in parallel
build_parallel() {
    local services=("$@")
    local pids=()
    local results=()
    
    print_info "Building ${#services[@]} services in parallel..."
    
    for service in "${services[@]}"; do
        (build_service "$service") &
        pids+=($!)
    done
    
    # Wait for all builds to complete
    local failed=0
    for i in "${!pids[@]}"; do
        local pid=${pids[$i]}
        local service=${services[$i]}
        
        if wait $pid; then
            results+=("$service: SUCCESS")
        else
            results+=("$service: FAILED")
            failed=$((failed + 1))
        fi
    done
    
    # Print results
    print_header "Parallel Build Results"
    for result in "${results[@]}"; do
        if [[ $result == *"SUCCESS"* ]]; then
            print_success "$result"
        else
            print_error "$result"
        fi
    done
    
    return $failed
}

# Main build logic
print_header "Build Configuration"
print_info "Force build: $FORCE_BUILD"
print_info "Specific service: ${SPECIFIC_SERVICE:-"all"}"
print_info "No cache: $NO_CACHE"
print_info "Parallel build: $PARALLEL_BUILD"
print_info "Registry: ${REGISTRY:-"local only"}"
print_info "Tag: $TAG"
print_info "Check only: $CHECK_ONLY"

# Determine which services to build
services_to_build=()
if [ -n "$SPECIFIC_SERVICE" ]; then
    if [[ -v SERVICES["$SPECIFIC_SERVICE"] ]]; then
        services_to_build=("$SPECIFIC_SERVICE")
    else
        print_error "Unknown service: $SPECIFIC_SERVICE"
        print_info "Available services: ${!SERVICES[*]}"
        exit 1
    fi
else
    services_to_build=($(printf '%s\n' "${!SERVICES[@]}" | sort))
fi

print_header "Services to Build"
for service in "${services_to_build[@]}"; do
    if needs_rebuild "$service" || [ "$FORCE_BUILD" = true ]; then
        if [ "$CHECK_ONLY" = true ]; then
            print_warning "$service (changes detected)"
        else
            print_info "$service (will build)"
        fi
    else
        print_success "$service (up to date)"
    fi
done

if [ "$CHECK_ONLY" = true ]; then
    print_success "Change detection completed"
    exit 0
fi

# Filter services that actually need building
services_needing_build=()
for service in "${services_to_build[@]}"; do
    if needs_rebuild "$service" || [ "$FORCE_BUILD" = true ]; then
        services_needing_build+=("$service")
    fi
done

if [ ${#services_needing_build[@]} -eq 0 ]; then
    print_success "All services are up to date!"
    exit 0
fi

print_header "Building ${#services_needing_build[@]} service(s)"

# Build services
if [ "$PARALLEL_BUILD" = true ] && [ ${#services_needing_build[@]} -gt 1 ]; then
    if build_parallel "${services_needing_build[@]}"; then
        print_success "All builds completed successfully!"
    else
        print_error "Some builds failed"
        exit 1
    fi
else
    # Sequential build
    failed_builds=()
    for service in "${services_needing_build[@]}"; do
        if ! build_service "$service"; then
            failed_builds+=("$service")
        fi
    done
    
    if [ ${#failed_builds[@]} -eq 0 ]; then
        print_success "All builds completed successfully!"
    else
        print_error "Failed builds: ${failed_builds[*]}"
        exit 1
    fi
fi

# Summary
print_header "Build Summary"
print_info "Build ID: $BUILD_ID"
print_info "Services built: ${#services_needing_build[@]}"
if [ -n "$REGISTRY" ]; then
    print_info "Images pushed to: $REGISTRY"
fi
print_success "Build process completed successfully!"
