# Development Scripts

This directory contains scripts to help manage your development environment.

## Hydra OAuth2 Client Management

### Quick Setup (Recommended)

```bash
# Complete environment setup (cluster + OAuth2 clients)
./scripts/setup-dev-environment.sh
```

### Individual Scripts

#### Create Hydra Clients
```bash
# Create all OAuth2 clients
./scripts/create-hydra-clients.sh
```

#### Verify Clients
```bash
# Check if all clients exist
./scripts/verify-hydra-clients.sh
```

#### Show Client Credentials
```bash
# Display all client credentials for reference
./scripts/show-client-credentials.sh
```

---

# Hydra OAuth2 Clients Setup Script

This script automatically creates the necessary OAuth2 clients in Hydra for your development environment.

## What it creates

The script creates three OAuth2 clients:

### 1. Dashboard Application
- **Client ID**: `4b41cd38-43ed-4e3a-9a88-bd384af21732`
- **Purpose**: SvelteKit dashboard authentication
- **Redirect URI**: `http://localhost:5173/auth/callback/hydra`
- **Grant Types**: `authorization_code`, `refresh_token`
- **Scopes**: `openid`, `offline`, `hydra.openid`, `introspect`

### 2. Swagger UI Documentation
- **Client ID**: `d39beaaa-9c53-48e7-b82a-37ff52127473`
- **Purpose**: API documentation OAuth2 authentication
- **Redirect URI**: `http://127.0.0.1:8080/v1/docs/oauth2-redirect.html`
- **Grant Types**: `authorization_code`, `refresh_token`
- **Scopes**: `openid`, `offline`, `profile`, `email`, `offline_access`

### 3. Oathkeeper Proxy
- **Client ID**: `761506cc-511f-411c-bf31-752efd8063b3`
- **Purpose**: Token introspection for the API gateway
- **Grant Types**: `client_credentials`
- **Scopes**: `introspect`

## Prerequisites

- k3s cluster is running
- Hydra is deployed and ready
- kubectl is configured to access the cluster

## Usage

```bash
# From the root of your project
./scripts/create-hydra-clients.sh
```

## When to run

Run this script after:
- Creating a new k3s cluster
- Deleting and recreating Hydra
- Any time the OAuth2 clients are missing

## Integration

This script is designed to be run after your cluster setup. You can integrate it into your deployment workflow:

```bash
# Example: Run after cluster creation
./k3s/scripts/deploy-all.sh
./scripts/create-hydra-clients.sh
```

## Troubleshooting

### Script hangs at "Waiting for Hydra admin"
- Check if Hydra pod is running: `kubectl get pods -n naytife-auth`
- Check Hydra logs: `kubectl logs -n naytife-auth deploy/hydra`

### "Client already exists" warnings
- This is normal - the script will delete and recreate existing clients

### Permission errors
- Ensure the script is executable: `chmod +x scripts/create-hydra-clients.sh`
- Ensure kubectl has proper cluster access

## Security Note

These client secrets are for **development only**. In production:
- Use different, secure client secrets
- Store secrets in a secure secret management system
- Use proper RBAC and network policies
