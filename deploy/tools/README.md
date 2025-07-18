# Skaffold Development Environment

This directory contains Skaffold configurations for local development of the Naytife platform.

## Overview

Skaffold provides a seamless development workflow that:
- Watches for file changes and automatically rebuilds/redeploys
- Provides port forwarding for easy access to services
- Streams logs from all services
- Integrates with the existing Kustomize deployment setup

## Quick Start

### Prerequisites

1. **Install Skaffold**:
   ```bash
   brew install skaffold
   ```

2. **Ensure k3s is running**:
   ```bash
   ./k3s/scripts/create-cluster.sh
   ```

3. **Validate environment**:
   ```bash
   ./deploy/scripts/validate-environment.sh local
   ```

### Start Development Environment

```bash
# Start development with file watching and hot reload
./deploy/scripts/deploy-skaffold.sh dev

# Start in debug mode
./deploy/scripts/deploy-skaffold.sh debug

# Build and deploy once (no file watching)
./deploy/scripts/deploy-skaffold.sh run

# Fast build mode with increased concurrency
./deploy/scripts/deploy-skaffold.sh fast
```

### Access Services

Once deployed, services are available at:
- **Backend**: http://localhost:8000
- **Auth Handler**: http://localhost:3000
- **Oathkeeper**: http://localhost:4456
- **Hydra Admin**: http://localhost:4445
- **PostgreSQL**: localhost:5432

## File Structure

```
/
├── skaffold.yaml                    # Main Skaffold configuration
└── deploy/
    ├── tools/
    │   ├── skaffold-local.yaml     # Local environment overrides
    │   └── skaffold-dev.yaml       # Development-specific settings
    └── scripts/
        ├── deploy-skaffold.sh      # Main deployment script
        ├── skaffold-utils.sh       # Utility commands
        └── validate-environment.sh # Environment validation
```

## Configuration Files

### Main Configuration (`skaffold.yaml`)

The main Skaffold configuration defines:
- **Build artifacts** for all custom services
- **File sync patterns** for hot reloading
- **Deployment** using Kustomize local overlay
- **Port forwarding** for development access
- **Profiles** for different development modes

### Environment-Specific Configurations

#### Local Override (`deploy/tools/skaffold-local.yaml`)
- Local build settings (no registry push)
- Local deployment flags
- Test configurations

#### Development Settings (`deploy/tools/skaffold-dev.yaml`)
- Debug build arguments
- Extended port forwarding
- Enhanced logging

## Development Workflow

### File Watching and Hot Reload

Skaffold automatically watches for changes in:

#### **Backend Service**
- `**/*.go` - All Go source files
- `docs/**` - Documentation files  
- `internal/**` - Internal package files
- `cmd/**` - Command files
- `config/**` - Configuration files

#### **Auth Handler, Store Deployer, Template Registry**
- `**/*.go` - All Go source files
- `go.mod` - Go module definition
- `go.sum` - Go module checksums

When files change:
1. **Go source files**: File sync to running container (1-3 seconds)
2. **go.mod/go.sum**: File sync + potential module reload
3. **Dockerfile changes**: Complete rebuild (30-90 seconds)

### Profiles

#### `local` Profile
- Local Docker builds (no registry push)
- Uses local k3s cluster
- Optimized for development speed

#### `debug` Profile
- Enables debug build flags
- Additional debugging ports
- Extended logging

#### `fast` Profile
- Increased build concurrency
- Forced deployments
- Optimized for rapid iteration

## Utility Commands

Use the utility script for common operations:

```bash
# Build all images without deploying
./deploy/scripts/skaffold-utils.sh build

# Render Kubernetes manifests
./deploy/scripts/skaffold-utils.sh render

# Deploy using pre-built images
./deploy/scripts/skaffold-utils.sh deploy

# Check deployment status
./deploy/scripts/skaffold-utils.sh status

# View logs from all services
./deploy/scripts/skaffold-utils.sh logs

# View logs from specific service
./deploy/scripts/skaffold-utils.sh logs backend

# Start debug mode
./deploy/scripts/skaffold-utils.sh debug

# Validate configuration
./deploy/scripts/skaffold-utils.sh validate

# Clean up Docker images and cache
./deploy/scripts/skaffold-utils.sh clean

# Delete all deployed resources
./deploy/scripts/skaffold-utils.sh delete
```

## Integration with Existing Setup

### Kustomize Integration
- Uses existing `deploy/overlays/local` configuration
- Inherits all current patches and configurations
- Maintains namespace structure (`naytife-local`)
- Respects existing image naming conventions

### Image Naming
Matches your current overlay configuration:
- `naytife/backend`
- `naytife/auth-handler`
- `docker.io/library/store-deployer`
- `template-registry`

### Service Discovery
Services are accessible within the cluster using:
- `backend.naytife-local.svc.cluster.local:8000`
- `auth-handler.naytife-local.svc.cluster.local:3000`
- `postgres.naytife-local.svc.cluster.local:5432`

## Performance Optimization

### Build Performance
- **BuildKit** enabled for faster builds
- **Layer caching** for Go modules
- **Concurrent builds** in fast mode
- **File sync** for rapid development

### Expected Performance
- **Initial build**: 2-5 minutes (all services)
- **Incremental build**: 10-30 seconds (single service)
- **File sync**: 1-3 seconds (Go file changes)
- **Full rebuild**: 30-90 seconds (dependency changes)

## Troubleshooting

### Common Issues

#### Build Failures
```bash
# Check Docker daemon
docker version

# Clear Docker cache
docker system prune -f

# Validate Dockerfile syntax
./deploy/scripts/skaffold-utils.sh validate
```

#### Deployment Issues
```bash
# Check cluster connectivity
kubectl get nodes

# Verify namespace
kubectl get namespace naytife-local

# Check deployment status
./deploy/scripts/skaffold-utils.sh status
```

#### Port Forwarding Problems
```bash
# Check if ports are in use
lsof -i :8000

# Verify service endpoints
kubectl get endpoints -n naytife-local
```

#### File Sync Issues
```bash
# Check file permissions
ls -la backend/

# Manual sync test
skaffold dev --no-sync

# Verbose output
skaffold dev -v debug
```

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
# Verbose Skaffold output
skaffold dev -v debug --profile=local

# Skip build and deploy directly
skaffold dev --skip-build

# Force rebuild all images
skaffold dev --cache-artifacts=false
```

## Development Tips

### Efficient Development
1. Use `dev` mode for active development with file watching
2. Use `run` mode for one-time testing
3. Use `debug` mode when debugging issues
4. Use `fast` mode for rapid iteration

### File Organization
- Keep frequently changed files in sync patterns
- Use `.dockerignore` to exclude unnecessary files
- Organize code to minimize rebuild triggers

### Resource Management
- Monitor resource usage with `kubectl top nodes`
- Clean up unused images regularly
- Use `cleanup` flags in dev mode

## Migration from Existing k3s Scripts

### Advantages of Skaffold
- **Automated file watching** vs manual rebuilds
- **Integrated port forwarding** vs manual kubectl commands  
- **Log streaming** from all services
- **Faster builds** with intelligent caching
- **Hot reload** for rapid development

### Maintaining Compatibility
- All existing Kustomize configurations work unchanged
- Same namespace and service structure
- Compatible with existing scripts and monitoring
- Same access patterns and service discovery

### Migration Steps
1. Validate environment: `./deploy/scripts/validate-environment.sh local`
2. Start with Skaffold: `./deploy/scripts/deploy-skaffold.sh dev`
3. Test all services are accessible
4. Verify file watching works with a small change
5. Use existing monitoring and debugging tools

The Skaffold setup enhances the development experience while maintaining full compatibility with your existing infrastructure and workflows.
