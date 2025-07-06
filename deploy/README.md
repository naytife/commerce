# Kubernetes Deployment Configuration

This directory contains the modern Kubernetes deployment configuration using Kustomize and SOPS for the Naytife platform.

## 🏗️ Architecture Overview

This deployment system implements a GitOps-ready, secure, and scalable Kubernetes deployment workflow using:

- **Kustomize**: For configuration management and environment overlays
- **SOPS**: For encrypted secrets management with age encryption
- **Environment Isolation**: Separate configurations for local, staging, and production

## 📁 Directory Structure

```
deploy/
├── base/                          # Base Kustomize configurations
│   ├── auth-handler/             # Authentication handler service
│   ├── backend/                  # Main backend API service
│   ├── postgres/                 # PostgreSQL database
│   ├── redis/                    # Redis cache
│   ├── hydra/                    # OAuth2 server
│   ├── oathkeeper/               # API gateway
│   ├── template-registry/        # Template registry service
│   ├── store-deployer/           # Store deployment service
│   ├── namespaces/               # Namespace definitions
│   └── kustomization.yaml        # Base kustomization file
├── overlays/                     # Environment-specific overlays
│   ├── local/                    # Local development environment
│   ├── staging/                  # Staging environment
│   └── production/               # Production environment
├── secrets/                      # SOPS encrypted secrets
│   ├── local/                    # Local environment secrets
│   ├── staging/                  # Staging environment secrets
│   └── production/               # Production environment secrets
├── scripts/                      # Deployment and utility scripts
│   ├── deploy.sh                 # Main deployment script
│   └── validate-phase1.sh        # Phase 1 validation script
└── tools/                        # Development tools configuration
```

## 🚀 Phase 1 Implementation Status

### ✅ Phase 1.1: Foundation Setup
- [x] Deploy directory structure created
- [x] Base Kustomize configurations for all services
- [x] Migrated existing YAML manifests to Kustomize base resources

### ✅ Phase 1.2: SOPS Setup for Secrets Management
- [x] SOPS configured with age encryption
- [x] Encryption keys generated for all environments (local, staging, production)
- [x] All secrets encrypted and stored in environment-specific directories
- [x] Hardcoded secrets removed from YAML manifests

### ✅ Phase 1.3: Environment Overlays
- [x] Kustomize overlays created for local, staging, and production
- [x] Environment-specific configurations defined
- [x] Proper resource naming conventions implemented
- [x] Environment-specific ingress and services configured

## 🌍 Environment Configurations

### Local Environment (`overlays/local/`)
- **Purpose**: Development and testing
- **Namespace**: `naytife-local`
- **Resource Prefix**: `local-`
- **Features**:
  - ImagePullPolicy set to Never (for local images)
  - Debug logging enabled
  - NodePort services for external access
  - Minimal resource requirements

### Staging Environment (`overlays/staging/`)
- **Purpose**: Pre-production testing
- **Namespace**: `naytife-staging`
- **Resource Prefix**: `staging-`
- **Features**:
  - Moderate resource allocations
  - Info-level logging
  - Ingress with TLS (staging certificates)
  - 2 replicas for key services

### Production Environment (`overlays/production/`)
- **Purpose**: Live production workloads
- **Namespace**: `naytife-production`
- **Resource Prefix**: `prod-`
- **Features**:
  - High resource allocations
  - Warning-level logging
  - Production ingress with TLS
  - 3 replicas for high availability
  - Network policies for security
  - Health checks and monitoring

## 🔐 Secrets Management

Secrets are managed using SOPS with age encryption. Each environment has its own encryption key:

- **Local**: `age1pynp2nwc45zjy6a7ka3vxghqxhac5v2506tjj695rvxfwj2fcfgs77ly9l`
- **Staging**: `age13ynzgj8jc2ddqj8jdq84crs6ev4yf2sa3srj00dg4y6cfk2c5p8s5rdlu8`
- **Production**: `age1lygw3utcj5eguktcjt583e2gpcgu4m7shv2mj2cyn93z2nggpv9sua67hu`

### Secret Files per Environment:
- `backend-secret.yaml` - Database and Redis connection strings
- `auth-handler-secret.yaml` - Authentication service secrets
- `postgres-secret.yaml` - Database credentials
- `redis-secret.yaml` - Redis authentication
- `hydra-secret.yaml` - OAuth2 server configuration
- `oathkeeper-secret.yaml` - API gateway secrets
- `cloudflare-secrets.yaml` - CDN and DNS credentials

## 🛠️ Usage

### Prerequisites
- `kubectl` - Kubernetes command-line tool
- `kustomize` - Configuration management tool
- `sops` - Secrets management tool
- `age` - Encryption tool for SOPS

### Deployment Commands

#### Deploy to Local Environment
```bash
./scripts/deploy.sh local
```

#### Deploy to Staging Environment
```bash
./scripts/deploy.sh staging
```

#### Deploy to Production Environment
```bash
./scripts/deploy.sh production
```

#### Dry Run (Preview Changes)
```bash
./scripts/deploy.sh <environment> --dry-run
```

### Validation

Run the Phase 1 validation script to ensure everything is properly configured:

```bash
./scripts/validate-phase1.sh
```

### Manual Kustomize Operations

#### Build Configuration
```bash
cd overlays/<environment>
kustomize build .
```

#### Apply Configuration
```bash
cd overlays/<environment>
kustomize build . | kubectl apply -f -
```

## 🔧 Configuration Management

### Adding New Environment Variables
1. Update the appropriate `configMapGenerator` in the overlay's `kustomization.yaml`
2. Add environment-specific values for each environment

### Adding New Secrets
1. Create the secret in the appropriate environment directory under `secrets/`
2. Encrypt with SOPS: `sops -e -i secrets/<environment>/new-secret.yaml`
3. Add reference in the overlay's `kustomization.yaml`

### Adding New Services
1. Create base configuration in `base/new-service/`
2. Add reference to `base/kustomization.yaml`
3. Create environment-specific patches in each overlay as needed

### Modifying Resource Allocations
Update the patch files in each environment overlay:
- `*-local-patch.yaml` for local environment
- `*-staging-patch.yaml` for staging environment  
- `*-production-patch.yaml` for production environment

## 🏷️ Resource Naming Convention

Resources are prefixed based on environment:
- **Local**: `local-<resource-name>`
- **Staging**: `staging-<resource-name>`
- **Production**: `prod-<resource-name>`

Namespaces follow the pattern:
- **Local**: `naytife-local`
- **Staging**: `naytife-staging`
- **Production**: `naytife-production`

## 🔍 Troubleshooting

### Common Issues

#### 1. SOPS Decryption Fails
```bash
# Check if age key is properly configured
sops -d secrets/local/backend-secret.yaml
```

#### 2. Kustomize Build Fails
```bash
# Validate kustomization.yaml syntax
cd overlays/<environment>
kustomize build . --validate
```

#### 3. Resource Not Found
```bash
# Check if all referenced files exist
find . -name "*.yaml" -exec ls -la {} \;
```

### Validation Commands

```bash
# Validate all environments
for env in local staging production; do
  echo "Validating $env..."
  cd overlays/$env && kustomize build . > /dev/null && echo "✅ $env OK"
done

# Check secret encryption
find secrets/ -name "*.yaml" -exec grep -L "ENC\[" {} \;
```

## 📊 Resource Overview

### Services Deployed
- **backend**: Main API service (Go)
- **authentication-handler**: Auth service (Go)
- **postgres**: Database (PostgreSQL 15)
- **redis**: Cache and session store
- **hydra**: OAuth2/OIDC server
- **oathkeeper**: API gateway and proxy
- **template-registry**: Template management service
- **store-deployer**: Store deployment automation

### Infrastructure Components
- **Ingress**: Environment-specific routing and TLS
- **NetworkPolicies**: Security controls (production only)
- **PersistentVolumes**: Data persistence (production)
- **ConfigMaps**: Environment configuration
- **Secrets**: Encrypted sensitive data

## 🔒 Security Features

### Production Security
- Network policies restrict inter-pod communication
- TLS termination at ingress level
- Encrypted secrets with environment-specific keys
- Resource quotas and limits
- Health checks and liveness probes

### Development Security
- Isolated namespaces per environment
- Separate encryption keys per environment
- No secrets in version control
- Audit trail for secret access

## 📈 Next Steps (Phase 2)

The next phase will include:
- Skaffold integration for local development
- Automated development scripts
- Integration testing framework
- Hot reloading and file watching

## 🆘 Support

For issues related to this deployment configuration:
1. Run the validation script: `./scripts/validate-phase1.sh`
2. Check logs: `kubectl logs -n <namespace> <pod-name>`
3. Verify configurations: `kustomize build overlays/<environment>`

## 📚 References

- [Kustomize Documentation](https://kustomize.io/)
- [SOPS Documentation](https://github.com/mozilla/sops)
- [Age Encryption](https://age-encryption.org/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
