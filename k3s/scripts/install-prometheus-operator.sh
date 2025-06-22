#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📊 Installing Prometheus Operator for Monitoring${NC}"
echo "=================================================="

# Check if helm is installed
if ! command -v helm &> /dev/null; then
    echo -e "${RED}❌ Helm is not installed. Please install Helm first.${NC}"
    echo "Visit: https://helm.sh/docs/intro/install/"
    exit 1
fi

echo -e "${YELLOW}📦 Adding Prometheus Operator Helm repository...${NC}"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

echo -e "${YELLOW}🔧 Creating monitoring namespace...${NC}"
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
kubectl label namespace monitoring name=monitoring --overwrite

echo -e "${YELLOW}🚀 Installing Prometheus Operator...${NC}"
helm install prometheus-operator prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.ruleSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.retention=7d \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=10Gi \
  --set grafana.adminPassword=admin123 \
  --set grafana.persistence.enabled=true \
  --set grafana.persistence.size=5Gi

echo -e "${YELLOW}⏳ Waiting for Prometheus Operator to be ready...${NC}"
kubectl wait --for=condition=available deployment/prometheus-operator-kube-p-operator -n monitoring --timeout=300s
kubectl wait --for=condition=available deployment/prometheus-operator-grafana -n monitoring --timeout=300s

echo -e "${GREEN}✅ Prometheus Operator installed successfully!${NC}"

echo -e "\n${BLUE}🔗 Access Points:${NC}"
echo "================================"
echo "  📊 Prometheus: kubectl port-forward -n monitoring svc/prometheus-operator-kube-p-prometheus 9090:9090"
echo "  📈 Grafana:    kubectl port-forward -n monitoring svc/prometheus-operator-grafana 3000:80"
echo "  🚨 AlertManager: kubectl port-forward -n monitoring svc/prometheus-operator-kube-p-alertmanager 9093:9093"

echo -e "\n${BLUE}📋 Grafana Credentials:${NC}"
echo "  Username: admin"
echo "  Password: admin123"

echo -e "\n${YELLOW}📝 Next Steps:${NC}"
echo "1. Enable ServiceMonitor in cloud-build deployment:"
echo "   - Uncomment the ServiceMonitor section in k3s/manifests/07-cloud-build/cloud-build.yaml"
echo "   - Apply: kubectl apply -f k3s/manifests/07-cloud-build/"
echo "2. Access Grafana and import cloud-build dashboards"
echo "3. Configure alerting rules for production monitoring"

echo -e "\n${GREEN}🎉 Monitoring infrastructure ready!${NC}"
