# Deployment Scripts Documentation

This directory contains enhanced deployment and development scripts for the Naytife platform using Kustomize-based deployments.

## Available Scripts

### Core Deployment Scripts

#### `deploy.sh`
Main deployment script with Kustomize and SOPS integration.

```bash
# Deploy to local environment
./deploy.sh local

# Deploy to staging with dry-run
./deploy.sh staging --dry-run

# Deploy to production
./deploy.sh production
```

**Features:**
- Kustomize-based configuration building
- SOPS secrets decryption
- Environment validation
- Deployment health checks
- Rollback support

#### `cleanup.sh`
Clean up deployed resources from any environment.

```bash
# Clean up local environment
./cleanup.sh local

# Preview what would be deleted
./cleanup.sh staging --dry-run

# Force cleanup without confirmation
./cleanup.sh production --force
```

**Features:**
- Safe resource deletion order
- Environment-specific safety checks
- Dry-run capability
- Graceful pod termination

### Development and Debugging Scripts

#### `status.sh`
Enhanced status monitoring for deployed services.

```bash
# Check local environment status
./status.sh local

# Detailed status for staging
./status.sh staging --detailed

# Quick status check
./status.sh
```

**Features:**
- Pod and deployment status
- Service health checks
- Resource usage monitoring
- Recent events display

#### `logs.sh`
Advanced log aggregation and viewing.

```bash
# View all service logs
./logs.sh

# View backend logs in staging
./logs.sh backend staging

# Follow logs with tail
./logs.sh backend local --follow --tail=50
```

**Features:**
- Service-specific log viewing
- Multi-pod log aggregation
- Follow mode for real-time logs
- Environment-aware log access

#### `debug.sh`
Comprehensive debugging and troubleshooting tool.

```bash
# Debug all services in local
./debug.sh local

# Debug specific service with details
./debug.sh local backend --detailed

# Debug production environment
./debug.sh production
```

**Features:**
- Deployment and pod analysis
- Event and condition checking
- Resource utilization review
- Problem identification

### Validation and Testing Scripts

#### `validate-environment.sh`
Environment configuration validation.

```bash
# Validate local environment
./validate-environment.sh local

# Validate staging setup
./validate-environment.sh staging
```

**Features:**
- Prerequisites checking
- Configuration validation
- Service health verification
- Comprehensive test suite

#### `test-integration.sh`
Integration testing suite for deployed services.

```bash
# Run all integration tests
./test-integration.sh local

# Test specific service
./test-integration.sh local --service=backend

# Verbose testing output
./test-integration.sh staging --verbose
```

**Features:**
- Service connectivity testing
- Database connectivity checks
- Authentication flow validation
- End-to-end workflow testing

#### `benchmark.sh`
Performance benchmarking tool.

```bash
# Benchmark all services
./benchmark.sh local

# Benchmark specific service
./benchmark.sh local --service=backend --duration=60s

# High concurrency test
./benchmark.sh staging --concurrency=50
```

**Features:**
- HTTP load testing with hey
- Response time analysis
- Resource usage monitoring
- Detailed performance reports

### Security and Secret Management

#### `sops-helper.sh`
Comprehensive SOPS encryption and secret management utilities.

```bash
# Edit encrypted secrets
./sops-helper.sh edit local postgres-secret.yaml

# View decrypted secrets
./sops-helper.sh view staging backend-secret.yaml

# Validate all encrypted files
./sops-helper.sh validate

# Generate new age key
./sops-helper.sh keygen

# Show public key
./sops-helper.sh keyshow

# Check key security
./sops-helper.sh keycheck
```

## Script Dependencies

### Required Tools
- `kubectl` - Kubernetes CLI
- `kustomize` - Configuration management
- `sops` - Secrets encryption
- `curl` - HTTP testing
- `hey` - Load testing (auto-installed via go)

### Optional Tools
- `go` - For installing additional tools
- Metrics server (for resource usage data)

## Environment Configuration

### Namespace Mapping
- **local**: `naytife`
- **staging**: `naytife-staging`
- **production**: `naytife-production`

### Service Mappings
| Service | Port | Health Endpoint |
|---------|------|-----------------|
| backend | 8000 | `/health` |
| auth-handler | 8080 | `/health` |
| hydra | 4444 | `/health/ready` |
| oathkeeper | 4456 | `/health/ready` |
| postgres | 5432 | N/A |
| redis | 6379 | N/A |
| store-deployer | 8090 | `/health` |
| template-registry | 8091 | `/health` |

## Usage Patterns

### Development Workflow

1. **Deploy to local environment:**
   ```bash
   ./deploy.sh local
   ```

2. **Check deployment status:**
   ```bash
   ./status.sh local --detailed
   ```

3. **Run integration tests:**
   ```bash
   ./test-integration.sh local
   ```

4. **Debug any issues:**
   ```bash
   ./debug.sh local
   ```

5. **View logs for troubleshooting:**
   ```bash
   ./logs.sh backend local --follow
   ```

### Performance Testing

1. **Run baseline benchmarks:**
   ```bash
   ./benchmark.sh local --duration=30s
   ```

2. **Test under load:**
   ```bash
   ./benchmark.sh local --concurrency=20 --duration=60s
   ```

3. **Check resource usage:**
   ```bash
   ./status.sh local --detailed
   ```

### Environment Validation

1. **Validate configuration:**
   ```bash
   ./validate-environment.sh staging
   ```

2. **Test service connectivity:**
   ```bash
   ./test-integration.sh staging --verbose
   ```

3. **Monitor deployment:**
   ```bash
   ./status.sh staging
   ```

## Error Handling

All scripts include comprehensive error handling with:
- **Exit codes**: 0 for success, non-zero for failures
- **Colored output**: Visual indicators for status
- **Verbose modes**: Detailed debugging information
- **Dry-run options**: Preview changes before execution

## Security Considerations

### Secrets Management
- All secrets are encrypted with SOPS
- Environment-specific encryption keys
- Automatic decryption during deployment
- No secrets stored in plain text

### Production Safety
- Confirmation prompts for production operations
- Force flags for automation
- Detailed logging of all operations
- Rollback capabilities

## Output and Logging

### Log Locations
- Benchmark results: `/tmp/naytife-benchmarks/`
- Temporary files: `/tmp/kustomize-output-*`
- Script logs: Console output with colored formatting

### Report Formats
- **Status**: Tabular console output
- **Benchmarks**: Markdown reports with metrics
- **Tests**: Pass/fail with detailed error messages
- **Debugging**: Structured troubleshooting information

## Troubleshooting

### Common Issues

1. **kubectl not connected:**
   ```bash
   kubectl config current-context
   kubectl cluster-info
   ```

2. **SOPS decryption fails:**
   ```bash
   ./sops-helper.sh check-keys
   sops --version
   ```

3. **Services not accessible:**
   ```bash
   ./debug.sh local
   ./validate-environment.sh local
   ```

4. **Performance issues:**
   ```bash
   ./benchmark.sh local
   ./status.sh local --detailed
   ```

### Getting Help

Each script includes a `--help` flag:
```bash
./deploy.sh --help
./status.sh --help
./test-integration.sh --help
```

## Migration from Legacy Scripts

The legacy k3s scripts in `/k3s/scripts/` provide similar functionality but use direct YAML manifests. The new scripts in `/deploy/scripts/` offer:

- **Enhanced security** with SOPS encryption
- **Environment management** with Kustomize overlays
- **Better debugging** with comprehensive tooling
- **Automated testing** with integration suites
- **Performance monitoring** with benchmarking tools

To migrate:
1. Use the new deployment scripts for all environments
2. Validate functionality with integration tests
3. Gradually phase out legacy scripts
4. Update CI/CD pipelines to use new scripts

## Contributing

When adding new scripts:
1. Follow the established patterns for colored output
2. Include comprehensive error handling
3. Add help documentation with `--help` flag
4. Ensure scripts are idempotent where possible
5. Add appropriate tests and validation
6. Update this documentation
