# Bucket Separation Migration Guide

## Overview

This guide documents the migration from a single `naytife-shops-static` bucket to separate `templates` and `stores` buckets for better organization and access control.

## Changes Made

### Before (Single Bucket Structure)
```
naytife-shops-static/
├── templates/
│   └── template_name/
│       └── version/
│           ├── manifest.json
│           └── assets/
└── stores/
    └── subdomain/
        ├── index.html
        └── data/
            └── metadata.json
```

### After (Separate Buckets)
```
templates/
└── template_name/
    └── version/
        ├── manifest.json
        └── assets/

stores/
├── subdomain/
│   ├── index.html
│   └── data/
│       └── metadata.json
├── products/
│   └── productId/
│       └── images/
└── shops/
    └── shopId/
        └── images/
```

## Services Updated

### 1. Template Registry Service
- **Bucket**: `templates`
- **Environment Variable**: `CLOUDFLARE_R2_BUCKET_NAME=templates`
- **Path Changes**: Files previously at `naytife-shops-static/templates/...` are now at `templates/...`
- **Kubernetes Secret Key**: `templates-bucket-name`

### 2. Store Deployer Service
- **Bucket**: `stores` 
- **Environment Variable**: `CLOUDFLARE_R2_BUCKET_NAME=stores`
- **Path Changes**: Files previously at `naytife-shops-static/stores/...` are now at `stores/...`
- **Kubernetes Secret Key**: `stores-bucket-name`

### 3. Dashboard Service
- **Bucket**: `stores` (for product and shop images)
- **Environment Variable**: `CLOUDFLARE_R2_BUCKET=stores`
- **Note**: Dashboard is not part of the Kubernetes cluster, environment needs to be updated separately

## Migration Steps

### Step 1: Create New Buckets
Create two new R2 buckets in your Cloudflare account:
1. `templates`
2. `stores`

### Step 2: Update Cloudflare R2 Configuration
The Kubernetes secrets have been updated to use separate bucket names:
```yaml
# k3s/manifests/08-template-system/cloudflare-secrets.yaml
data:
  templates-bucket-name: dGVtcGxhdGVz  # base64: templates
  stores-bucket-name: c3RvcmVz        # base64: stores
```

### Step 3: Migrate Existing Data
Run the following commands to migrate data from the old bucket structure:

#### Migrate Templates
```bash
# Copy all templates from old bucket to new templates bucket
rclone copy r2:naytife-shops-static/templates/ r2:templates/ --progress

# Verify the copy
rclone ls r2:templates/
```

#### Migrate Stores
```bash
# Copy all stores from old bucket to new stores bucket
rclone copy r2:naytife-shops-static/stores/ r2:stores/ --progress

# Copy any existing product images
rclone copy r2:naytife-shops-static/products/ r2:stores/products/ --progress

# Copy any existing shop images
rclone copy r2:naytife-shops-static/shops/ r2:stores/shops/ --progress

# Verify the copy
rclone ls r2:stores/
```

### Step 4: Update Dashboard Environment Variables
Update your dashboard application's environment variables:
```bash
# Old configuration (single bucket)
CLOUDFLARE_R2_BUCKET=naytife-shops-static

# New configuration (separate buckets)
CLOUDFLARE_R2_TEMPLATES_BUCKET=templates
CLOUDFLARE_R2_STORES_BUCKET=stores
```

**Note**: The dashboard API endpoints now use `CLOUDFLARE_R2_STORES_BUCKET` for product images and shop assets.

### Step 5: Deploy Updated Services
```bash
# Apply the updated Kubernetes manifests
kubectl apply -f k3s/manifests/08-template-system/

# Restart the services to pick up new environment variables
kubectl rollout restart deployment/template-registry -n naytife
kubectl rollout restart deployment/store-deployer -n naytife

# Verify the services are using the correct buckets
kubectl logs deployment/template-registry -n naytife | grep "bucket"
kubectl logs deployment/store-deployer -n naytife | grep "bucket"
```

### Step 6: Update Public URLs (if needed)
If you have public URLs pointing to the old bucket structure, you'll need to update them:

- Template assets: `https://pub-xxx.r2.dev/templates/...` → `https://pub-xxx.r2.dev/...`
- Store assets: `https://pub-xxx.r2.dev/stores/...` → `https://pub-xxx.r2.dev/...`

### Step 7: Test the Migration
1. Upload a new template through the template registry
2. Deploy a store through the store deployer
3. Upload product/shop images through the dashboard
4. Verify all files are going to the correct buckets

### Step 8: Cleanup (Optional)
Once you've verified everything works correctly, you can remove the old data:
```bash
# Remove templates from old bucket
rclone delete r2:naytife-shops-static/templates/ --dry-run
# rclone delete r2:naytife-shops-static/templates/  # Remove --dry-run when ready

# Remove stores from old bucket
rclone delete r2:naytife-shops-static/stores/ --dry-run
# rclone delete r2:naytife-shops-static/stores/    # Remove --dry-run when ready
```

## Benefits

1. **Better Organization**: Templates and stores are logically separated
2. **Access Control**: Different buckets can have different access policies
3. **Billing Separation**: Easier to track costs per service type
4. **Scalability**: Independent scaling and management of template vs store assets
5. **Security**: Reduce blast radius if one bucket is compromised

## Rollback Plan

If you need to rollback:

1. Update the Kubernetes secrets back to use `bucket-name: bmF5dGlmZS1zaG9wcy1zdGF0aWM=`
2. Update both services to use `bucket-name` instead of the separate keys  
3. Update dashboard environment variables:
   ```bash
   # Remove
   CLOUDFLARE_R2_TEMPLATES_BUCKET=templates
   CLOUDFLARE_R2_STORES_BUCKET=stores
   
   # Add back
   CLOUDFLARE_R2_BUCKET=naytife-shops-static
   ```
4. Restart the services

## Troubleshooting

### Service can't connect to bucket
- Check that the bucket names are correctly base64 encoded in the secrets
- Verify the bucket exists in your Cloudflare R2 account
- Check service logs: `kubectl logs deployment/[service-name] -n naytife`

### Files not found after migration
- Verify the rclone copy completed successfully
- Check that file paths are correct (no extra `/templates/` or `/stores/` prefix)
- Test with a simple file listing: `rclone ls r2:[bucket-name]/`

### Dashboard uploads failing
- Check the dashboard environment variables
- Verify the stores bucket exists and has proper permissions
- Check browser console for detailed error messages
