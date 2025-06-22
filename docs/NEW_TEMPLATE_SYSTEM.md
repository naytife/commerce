# New Template System Architecture

## Migration Status: ✅ COMPLETED

**Migration completed successfully on June 20, 2025**

### Migration Results:
- ✅ All 4 old services removed (template-builder, asset-manager, data-updater, old store-deployer)
- ✅ New template-registry service deployed and running (1 replica)
- ✅ Updated store-deployer service deployed and running (2 replicas)
- ✅ Template upload functionality tested and working
- ✅ Template discovery and listing working
- ✅ End-to-end store deployment tested successfully
- ✅ Test store deployed: `test-migration.naytife.com` using template_1 v1.0.0

### Performance Metrics:
- Template upload: < 5 seconds for ~79 assets (927KB total)
- Store deployment: ~1m52s for full asset copying and deployment
- Services: template-registry (21MB image), store-deployer (10.7MB image)
- Memory usage: template-registry (256Mi limit), store-deployer (512Mi limit)

---

## Overview

The template system has been redesigned from a 4-service architecture to a simplified 2-service architecture that is more maintainable, efficient, and production-ready.

## Architecture Changes

### Old Architecture (Removed)
- `template-builder` - Built templates (now templates are pre-built locally)
- `asset-manager` - Managed asset uploads and deduplication  
- `data-updater` - Updated store data separately
- `store-deployer` - Deployed stores with limited functionality

### New Architecture (Current)
1. **Template Registry Service** (`template-registry`)
   - Manages pre-built template assets and metadata
   - Handles template uploads via API
   - Provides template versioning and discovery
   - Stores templates in R2 at `templates/{template_name}/{version}/`

2. **Store Deployer Service** (`store-deployer`) - Enhanced
   - Deploys complete stores by copying templates and generating data
   - Handles both full deployments and data-only updates
   - Stores deployed sites in R2 at `stores/{shop_subdomain}/`
   - Integrates with template registry for template discovery

## R2 Storage Structure

```
bucket/
├── templates/
│   └── {template_name}/
│       ├── latest                    # Pointer to latest version
│       └── {version}/
│           ├── manifest.json         # Template metadata and asset list
│           ├── index.html           # Static template files
│           ├── _app/                # Svelte app files
│           └── ...                  # Other template assets
└── stores/
    └── {shop_subdomain}/
        ├── index.html               # Copied from template
        ├── _app/                    # Copied template assets
        ├── ...                      # Other static files
        └── data/                    # Store-specific data
            ├── shop.json            # Shop information
            ├── products.json        # Product catalog
            ├── settings.json        # Store configuration
            └── metadata.json        # Deployment metadata
```

## API Endpoints

### Template Registry Service (Port 9001)

- `GET /templates` - List all available templates
- `GET /templates/{template_name}` - Get template details and versions
- `GET /templates/{template_name}/versions` - List template versions
- `GET /templates/{template_name}/latest` - Get latest template version
- `GET /templates/{template_name}/versions/{version}` - Get specific version
- `POST /templates/upload` - Upload new template version
- `GET /health` - Health check

### Store Deployer Service (Port 9003)

- `POST /deploy` - Deploy a new store
- `POST /redeploy/{subdomain}` - Redeploy existing store
- `POST /update-data/{subdomain}` - Update store data only
- `GET /status/{subdomain}` - Get deployment status
- `GET /templates` - List available templates (proxied from registry)
- `GET /health` - Health check

## Deployment Flow

### New Store Deployment
1. Backend calls store-deployer `/deploy` endpoint
2. Store-deployer contacts template-registry for latest version
3. Store-deployer retrieves template manifest from R2
4. Assets are copied from template location to store location in R2
5. Store data is fetched from backend and uploaded as JSON files
6. Deployment response includes status and metadata

### Data-Only Updates
1. Backend calls store-deployer `/update-data/{subdomain}` endpoint
2. Store-deployer fetches fresh data from backend
3. Only data files in `stores/{subdomain}/data/` are updated
4. Static assets remain unchanged (faster updates)

## Template Upload Process

### Manual Upload
```bash
# Using the upload script
./scripts/upload-template.sh template_1 v1.2.0 "New template version"

# Or direct curl
curl -X POST \
  -F "template_name=template_1" \
  -F "version=v1.2.0" \
  -F "description=New template version" \
  -F "assets=@template.tar.gz" \
  http://template-registry:9001/templates/upload
```

### Build & Upload Workflow
1. Make changes to template in `templates/{template_name}/`
2. Build template: `cd templates/{template_name} && npm run build`
3. Upload to registry: `./scripts/upload-template.sh {template_name}`
4. Template is now available for new deployments

## Environment Variables

### Template Registry
- `CLOUDFLARE_R2_ACCESS_KEY_ID` - R2 access key
- `CLOUDFLARE_R2_SECRET_ACCESS_KEY` - R2 secret key  
- `CLOUDFLARE_R2_ENDPOINT` - R2 endpoint URL
- `CLOUDFLARE_R2_BUCKET_NAME` - R2 bucket name
- `PORT` - Service port (default: 9001)

### Store Deployer
- `CLOUDFLARE_R2_ACCESS_KEY_ID` - R2 access key
- `CLOUDFLARE_R2_SECRET_ACCESS_KEY` - R2 secret key
- `CLOUDFLARE_R2_ENDPOINT` - R2 endpoint URL  
- `CLOUDFLARE_R2_BUCKET_NAME` - R2 bucket name
- `TEMPLATE_REGISTRY_URL` - Template registry URL (default: http://template-registry:9001)
- `BACKEND_URL` - Backend GraphQL URL (default: http://backend:8002)
- `PORT` - Service port (default: 9003)

## Migration Guide

### From Old to New Architecture

1. **Backup existing data** (if needed)
2. **Run migration script:**
   ```bash
   cd k3s/scripts
   ./migrate-to-new-template-system.sh
   ```
3. **Verify services are running:**
   ```bash
   kubectl get pods -n commerce -l component=template-system
   ```
4. **Test template upload:**
   ```bash
   kubectl port-forward -n commerce svc/template-registry 9001:9001 &
   ./scripts/upload-template.sh template_1
   ```

### Key Changes in Backend Integration

- Store deployer URL remains the same: `http://store-deployer:9003`
- API endpoints remain compatible
- No changes needed in backend deployment handlers
- Data updates are now faster (data-only updates vs full redeployment)

## Monitoring & Troubleshooting

### Health Checks
```bash
# Template Registry
curl http://template-registry:9001/health

# Store Deployer  
curl http://store-deployer:9003/health
```

### View Logs
```bash
# Template Registry logs
kubectl logs -n commerce -l app=template-registry -f

# Store Deployer logs
kubectl logs -n commerce -l app=store-deployer -f
```

### Common Issues

1. **Template not found**: Ensure template is uploaded to registry
2. **R2 access denied**: Check R2 credentials in secrets
3. **Service not responding**: Check pod status and resource limits
4. **Asset copy failures**: Verify template manifest and R2 permissions

## Performance Improvements

- **Faster deployments**: Direct R2-to-R2 copying instead of multi-service calls
- **Reduced complexity**: 2 services instead of 4
- **Better caching**: Proper cache headers for static assets vs data
- **Data-only updates**: Skip asset copying for data changes
- **Template reuse**: Templates stored once, referenced by multiple stores

## Security Considerations

- All services run with non-root users
- Read-only root filesystems where possible
- Minimal resource limits and security contexts
- R2 credentials stored in Kubernetes secrets
- Internal service communication only (no external exposure)
