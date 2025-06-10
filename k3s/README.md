# Naytife K3s Deployment

Clean, organized Kubernetes deployment for local development.

## Architecture

```
k3s/
├── scripts/           # Deployment and management scripts
├── manifests/         # Kubernetes manifests organized by service
│   ├── 00-namespaces/ # Namespaces (deployed first)
│   ├── 01-postgres/   # PostgreSQL database
│   ├── 02-redis/      # Redis cache/queue
│   ├── 03-hydra/      # OAuth2 server
│   ├── 04-oathkeeper/ # API Gateway
│   ├── 05-auth-handler/ # Authentication handler
│   ├── 06-backend/    # Backend API
│   └── 07-cloud-build/ # Site builder service
└── configs/           # Configuration files
```

## Services

- **PostgreSQL** (port 5432): Main database
- **Redis** (port 6379): Cache and queue storage
- **Hydra** (port 4444/4445): OAuth2 server
- **Oathkeeper** (port 8080): API Gateway and proxy
- **Auth Handler** (port 3000): Login/consent handler
- **Backend** (port 8000): Main API server
- **Cloud Build** (port 9000): Static site generator

## Quick Start

```bash
# Deploy all services
./scripts/deploy-all.sh

# Deploy individual services
./scripts/deploy-service.sh postgres
./scripts/deploy-service.sh redis
./scripts/deploy-service.sh hydra
./scripts/deploy-service.sh oathkeeper
./scripts/deploy-service.sh auth-handler
./scripts/deploy-service.sh backend
./scripts/deploy-service.sh cloud-build

# Check status
./scripts/status.sh

# Cleanup
./scripts/cleanup.sh
```

## Development Workflow

1. **Build images**: `./scripts/build-images.sh`
2. **Deploy services**: `./scripts/deploy-all.sh`
3. **Test deployment**: `./scripts/test-deployment.sh`
4. **View logs**: `./scripts/logs.sh [service-name]`

## Service URLs

- API Gateway: http://127.0.0.1:8080
- Auth Handler: http://127.0.0.1:3000
- Backend API: http://127.0.0.1:8000
- Cloud Build: http://127.0.0.1:9000
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- Hydra Public: localhost:4444
- Hydra Admin: localhost:4445
