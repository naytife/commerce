# Naytife K3s Deployment Overview

This setup provides a clean, organized Kubernetes deployment for local development and testing, using [k3d](https://k3d.io/) to run a lightweight Kubernetes cluster in Docker. The structure is modular, with each service and configuration clearly separated.

## Directory Structure

```
k3s/
├── README.md
├── configs/         # Configuration files (resource configs, secrets, etc.)
├── manifests/       # Kubernetes manifests, organized by service
│   ├── 00-namespaces/      # Namespace definitions
│   ├── 01-postgres/        # PostgreSQL database manifests
│   ├── 02-redis/           # Redis cache/queue manifests
│   ├── 03-hydra/           # OAuth2 server manifests
│   ├── 04-oathkeeper/      # API Gateway manifests
│   ├── 05-auth-handler/    # Authentication handler manifests
│   ├── 06-backend/         # Backend API manifests
│   └── 08-template-system/ # Template system services
└── scripts/         # Shell scripts for cluster and service management
```

---

## 1. `README.md`

- **Purpose:** High-level documentation and quickstart guide.
- **Contents:**  
  - Architecture diagram and directory explanation.
  - List of core services and their ports.
  - Quickstart commands for deploying all or individual services.

---

## 2. `configs/`

- **Purpose:** Store configuration files for the cluster and services.
- **Example:**  
  - `arm-optimized-resources.yaml` (not detailed here): Resource requests/limits for ARM-based cloud VMs.

---

## 3. `manifests/`

Organized by service, each subdirectory contains Kubernetes YAML files for deploying and configuring that service.

### 3.1. `00-namespaces/`

- **File:** `namespaces.yaml`
- **Purpose:** Defines all Kubernetes namespaces used by the platform.
- **Usage:** Deploy this first to ensure all other resources are created in the correct namespaces.

### 3.2. `01-postgres/`

- **File:** `postgres.yaml`
- **Purpose:** Deploys the PostgreSQL database, including StatefulSet, Service, and PersistentVolumeClaim.
- **Notes:**  
  - Exposes port 5432.
  - Uses persistent storage for data durability.

### 3.3. `02-redis/`

- **File:** `redis.yaml`
- **Purpose:** Deploys Redis for caching and queueing.
- **Notes:**  
  - Exposes port 6379.
  - Can be used by multiple services for caching or as a message broker.

### 3.4. `03-hydra/`

- **Files:** `hydra.yaml`, `hydra-clients.yaml`
- **Purpose:**  
  - `hydra.yaml`: Deploys the [ORY Hydra](https://www.ory.sh/hydra/) OAuth2 server (public/admin endpoints).
  - `hydra-clients.yaml`: Defines OAuth2 clients for the platform.
- **Notes:**  
  - Exposes ports 4444 (public) and 4445 (admin).
  - Used for authentication and authorization flows.

### 3.5. `04-oathkeeper/`

- **File:** `oathkeeper.yaml`
- **Purpose:** Deploys [ORY Oathkeeper](https://www.ory.sh/oathkeeper/) as an API gateway and reverse proxy.
- **Notes:**  
  - Handles authentication, authorization, and request routing.
  - Exposes port 8080.

### 3.6. `05-auth-handler/`

- **File:** `auth-handler.yaml`
- **Purpose:** Deploys a custom authentication handler for login and consent flows.
- **Notes:**  
  - Integrates with Hydra for OAuth2 flows.
  - Exposes port 3000.

### 3.7. `06-backend/`

- **Files:** `backend.yaml`, `backend-migration.yaml`
- **Purpose:**  
  - `backend.yaml`: Deploys the main API server.
  - `backend-migration.yaml`: Runs database migrations as a Kubernetes Job.
- **Notes:**  
  - Exposes port 8000.
  - Migration job ensures the database schema is up to date before the API starts.

### 3.8. `08-template-system/`

- **Files:**  
  - `cloudflare-secrets.yaml`: Stores secrets for Cloudflare integration.
  - `store-deployer.yaml`: Deploys the store deployer service.
  - `template-registry.yaml`: Deploys the template registry service.
- **Notes:**  
  - These services are part of a microservice-based template system for site generation and deployment.

---

## 4. `scripts/`

A suite of Bash scripts to automate cluster and service management. All scripts use color-coded output and provide helpful prompts.

### Key Scripts

- **`build-images.sh`**  
  Builds Docker images for all Naytife services, tagging them for local use.

- **`build-template-services.sh`**  
  Builds and packages the template system microservices (`template-registry`, `store-deployer`).

- **`cleanup.sh`**  
  Deletes the k3d cluster and all associated resources. Prompts for confirmation unless forced.

- **`create-cluster.sh`**  
  Creates a new k3d cluster with all required port mappings. If a cluster exists, prompts to delete and recreate.

- **`deploy-all.sh`**  
  Deploys all services in the correct order, waiting for each to become ready. Handles dependencies and readiness checks.

- **`deploy-complete.sh`**  
  Checks prerequisites, ensures environment files exist, and runs a full deployment workflow.

- **`install-prometheus-operator.sh`**  
  Installs the Prometheus Operator for monitoring using Helm. Sets up monitoring and Grafana dashboards.

- **`load-images.sh`**  
  Loads locally built Docker images into the k3d cluster for use by Kubernetes.

- **`logs.sh`**  
  Fetches logs for any service, with options to follow logs and select namespaces.

- **`manage-migrations.sh`**  
  Manages database migrations for the backend, including running, rolling back, and checking migration status.

- **`setup-port-forwards.sh`**  
  Sets up local port forwards for all major services, making them accessible on localhost.

- **`status.sh`**  
  Shows the status of the cluster, namespaces, pods, services, and deployments. Includes health checks.

---

## 5. Service Overview

| Service           | Namespace      | Port(s) | Description                                 |
|-------------------|---------------|---------|---------------------------------------------|
| PostgreSQL        | naytife       | 5432    | Main database                               |
| Redis             | naytife       | 6379    | Cache and queue                             |
| Hydra             | naytife-auth  | 4444/5  | OAuth2 server (public/admin)                |
| Oathkeeper        | naytife-auth  | 8080    | API Gateway and proxy                       |
| Auth Handler      | naytife-auth  | 3000    | Login/consent handler                       |
| Backend           | naytife       | 8000    | Main API server                             |
| Store Deployer    | naytife       | 9003    | Deploys store instances                     |
| Template Registry | naytife       | 9001    | Manages template storage and retrieval      |

---

## 6. Deployment Workflow

1. **Build Docker Images:**  
   Run `./scripts/build-images.sh` and `./scripts/build-template-services.sh` to build all service images.

2. **Create Cluster:**  
   Run `./scripts/create-cluster.sh` to create a new k3d cluster with all necessary port mappings.

3. **Load Images:**  
   Run `./scripts/load-images.sh` to load the built images into the cluster.

4. **Deploy Services:**  
   Run `./scripts/deploy-all.sh` to deploy all services in order.

5. **Set Up Port Forwards:**  
   Run `./scripts/setup-port-forwards.sh` to access services on localhost.

6. **Monitor and Manage:**  
   - Use `./scripts/status.sh` to check cluster and service status.
   - Use `./scripts/logs.sh` to view logs.
   - Use `./scripts/manage-migrations.sh` to manage database migrations.
   - Use `./scripts/install-prometheus-operator.sh` to set up monitoring.

7. **Cleanup:**  
   Run `./scripts/cleanup.sh` to tear down the cluster and clean up resources.

---

## 7. Best Practices & Tips

- **Namespace First:** Always deploy namespaces before other resources.
- **Secrets:** Store sensitive data (e.g., Cloudflare credentials) in separate secrets files and do not commit them to version control.
- **Resource Management:** Adjust resource requests/limits in `configs/` for your hardware/cloud environment.
- **Monitoring:** Use Prometheus and Grafana for observability.
- **Modularity:** Each service is isolated in its own manifest directory for clarity and maintainability.

---

## 8. Replicating This Setup

- Install prerequisites: `k3d`, `kubectl`, `docker`, `helm`.
- Clone the repository and review the `README.md` for quickstart instructions.
- Follow the deployment workflow above.
- Adjust configuration and manifests as needed for your environment (e.g., ports, resource limits, secrets).

---

**This documentation should provide a comprehensive understanding of your `k3s` setup and enable others to create a similar environment.** If you need more details on any specific file or service, let me know!
