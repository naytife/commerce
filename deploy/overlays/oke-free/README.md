# OKE Always Free Deployment - Quick Start Guide

## 🚀 Deploy Your Application to Oracle Cloud (Always Free)

This guide will help you deploy your Naytife platform to Oracle Cloud Infrastructure using only Always Free services.

## ✅ What You Get (All FREE)
- **Full Kubernetes Cluster** (OKE - Oracle Kubernetes Engine)
- **ARM Compute** (4 cores, 24GB RAM total)
- **Load Balancer** (10 Mbps)
- **Block Storage** (200GB for databases)
- **Container Registry** (OCIR)
- **Your PostgreSQL Database** (containerized, just like local)

## 📋 Prerequisites

1. **OCI Account** with Always Free tier
2. **OCI CLI** installed and configured
3. **Docker** with buildx support
4. **kubectl** and **kustomize** installed

## 🏗️ Step 1: Setup OKE Cluster

### Create OKE Cluster in OCI Console
1. Go to **Developer Services → Kubernetes Clusters (OKE)**
2. Click **Create Cluster**
3. Choose **Quick Create**
4. Select **VM.Standard.A1.Flex** (ARM - Always Free)
5. Set node pool: **2 nodes, 2 OCPU, 12GB RAM each** (or 1 node with 4 OCPU, 24GB)
6. Create the cluster

### Setup kubectl Access
```bash
# Create kubeconfig (replace <cluster-id> and <region>)
oci ce cluster create-kubeconfig \
  --cluster-id <cluster-id> \
  --file $HOME/.kube/config-oke \
  --region <region> \
  --token-version 2.0.0

# Verify connection
export KUBECONFIG=$HOME/.kube/config-oke
kubectl get nodes
```

## 🔧 Step 2: Configure Your Environment

### Get Your OCI Details
```bash
# Get tenancy namespace
oci os ns get

# Get your region
oci iam region list
```

### Update Configuration Files
1. Edit `deploy/overlays/oke-free/kustomization.yaml`
2. Replace `<region>` and `<tenancy>` with your actual values:
   ```yaml
   images:
     - name: naytife/backend
       newName: us-ashburn-1.ocir.io/mytenancy/naytife/backend
       newTag: arm64-latest
   ```

### Setup OCIR Access
```bash
# Login to Oracle Container Registry
docker login <region>.ocir.io
# Username: <tenancy-namespace>/<your-username>
# Password: <your-auth-token>
```

## 🐳 Step 3: Build and Deploy

### Run Setup Check
```bash
cd deploy/scripts
./oke-setup.sh
```

### Build and Push ARM64 Images
```bash
# Set your environment variables
export OCI_REGION="us-ashburn-1"      # Your region
export OCI_TENANCY="mytenancy"        # Your tenancy namespace

# Build and push all images
./oke-build-push.sh
```

### Deploy to OKE
```bash
./oke-deploy.sh
```

## 🎯 Step 4: Access Your Application

### Get Load Balancer IP
```bash
kubectl get ingress naytife-ingress -n naytife-oke-free
```

### Add to /etc/hosts (for testing)
```bash
# Add these lines to /etc/hosts
<LOAD_BALANCER_IP> api.naytife-oke.dev
<LOAD_BALANCER_IP> auth.naytife-oke.dev
<LOAD_BALANCER_IP> oauth.naytife-oke.dev
<LOAD_BALANCER_IP> gateway.naytife-oke.dev
```

### Test Your API
```bash
curl http://api.naytife-oke.dev/health
```

## 📊 Monitoring and Management

### Check Pod Status
```bash
kubectl get pods -n naytife-oke-free -o wide
```

### View Logs
```bash
kubectl logs -f deployment/backend -n naytife-oke-free
```

### Scale Services
```bash
kubectl scale deployment backend --replicas=2 -n naytife-oke-free
```

### Resource Usage
```bash
kubectl top nodes
kubectl top pods -n naytife-oke-free
```

## 🔄 Updates and Maintenance

### Update Application
1. Make code changes
2. Build new images: `./oke-build-push.sh`
3. Redeploy: `./oke-deploy.sh`

### Update Configuration
1. Edit overlay files in `deploy/overlays/oke-free/`
2. Apply changes: `kustomize build deploy/overlays/oke-free | kubectl apply -f -`

## 🧹 Cleanup

### Remove Application (Keep Data)
```bash
./oke-teardown.sh
# Choose to preserve persistent volumes
```

### Complete Cleanup
```bash
./oke-teardown.sh
# Choose to delete everything including data
```

## 💡 Tips for Success

### Resource Management
- **Start Small**: Begin with 1 replica of each service
- **Monitor Usage**: Use `kubectl top` to watch resource consumption
- **Scale Gradually**: Increase replicas only if needed

### Cost Control
- **Always Free Resources**: Everything in this setup is free
- **Monitor OCI Console**: Check for any unexpected charges
- **Set Billing Alerts**: Configure alerts for any overages

### Database Management
- **Backup Strategy**: Your PostgreSQL data is in persistent volumes
- **Performance**: ARM instances provide excellent database performance
- **Scaling**: Can easily add more storage if needed

## 🚀 Production Migration

When ready for production:
1. **Create Production OKE**: Use paid instances for more resources
2. **Copy Overlay**: Use `oke-free` as template for `oke-production`
3. **Add CI/CD**: GitHub Actions can deploy to OKE
4. **Enable SSL**: Add TLS certificates to ingress
5. **Monitoring**: Add Prometheus/Grafana for observability

## 📖 Additional Resources

- **OCI Documentation**: [Oracle Cloud Infrastructure Docs](https://docs.oracle.com/en-us/iaas/)
- **OKE Documentation**: [Oracle Kubernetes Engine Docs](https://docs.oracle.com/en-us/iaas/Content/ContEng/home.htm)
- **Always Free Tier**: [OCI Always Free Resources](https://www.oracle.com/cloud/free/)

## 🆘 Troubleshooting

### Common Issues
1. **Images not pulling**: Check OCIR login and image names
2. **Pods pending**: Check node resources and ARM compatibility
3. **Services not accessible**: Check ingress and load balancer status
4. **Database connection issues**: Verify service names and networking

### Get Help
```bash
# Describe problematic resources
kubectl describe pod <pod-name> -n naytife-oke-free

# Check events
kubectl get events -n naytife-oke-free --sort-by='.lastTimestamp'

# Check logs
kubectl logs deployment/<service-name> -n naytife-oke-free
```

---

**🎉 Congratulations!** You now have a production-grade Kubernetes cluster running your application on Oracle Cloud's Always Free tier!
