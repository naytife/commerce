
# ✅ Zalando Postgres Operator 10/10 Checklist (k3s + ARM + Oracle Cloud)

This checklist guides you to deploy a **production-grade PostgreSQL cluster** using the **Zalando Postgres Operator** on your **k3s cluster running on ARM VMs in Oracle Cloud (Ampere A1)**.

---

## 🔧 1. Cluster & Node Setup

- [ ] Use Oracle ARM Ampere A1 VMs (e.g., 4-core / 24GB RAM)
- [ ] Provision 2–3 nodes (for HA / PodAntiAffinity)
- [ ] Use OCI Block Volumes via CSI driver for persistent storage
- [ ] Ensure nodes support `linux/arm64`
- [ ] Enable persistent storage with Oracle CSI driver
- [ ] Install k3s with `--disable traefik` if managing your own ingress
- [ ] Confirm PVC + CSIDriver support is active
- [ ] Label nodes `postgres=true` if isolating DB workloads

---

## 🛠️ 2. Install Zalando Postgres Operator

- [ ] Clone [Zalando Operator GitHub](https://github.com/zalando/postgres-operator)
- [ ] Apply CRDs: `kubectl apply -f manifests/01-crds.yaml`
- [ ] Deploy operator: `kubectl apply -f manifests/operator.yaml`
- [ ] Create namespace (e.g., `database`)
- [ ] Use multi-arch image like `ghcr.io/zalando/spilo-15:3.1-p1` if needed
- [ ] Confirm operator pod is running

---

## 🧠 3. Deploy PostgreSQL Cluster

- [ ] Create `Postgresql` CRD with `numberOfInstances: 3`
- [ ] Set `enableMasterLoadBalancer: false`
- [ ] Define users, DBs, resources in the CRD
- [ ] Use appropriate storage class (e.g., `oci`)
- [ ] Use `teamId` to logically group DBs

---

## 💾 4. Persistent Storage

- [ ] Confirm OCI CSI plugin is working
- [ ] Ensure PVCs are bound and healthy
- [ ] Validate data persists after pod/node restarts

---

## 🔁 5. Backup + PITR with WAL-G

- [ ] Create OCI Object Storage bucket (e.g., `pg-backups`)
- [ ] Create Kubernetes Secret with OCI access credentials
- [ ] Configure WAL-G env vars in CRD or Operator config
- [ ] Set logical backup schedule annotations
- [ ] Validate backup + restore flow

---

## 🔐 6. Security & Secrets

- [ ] Use Kubernetes Secrets for credentials
- [ ] Enable TLS (internal + client-side)
- [ ] Integrate with `cert-manager` for TLS automation if needed
- [ ] Restrict RBAC access to Secrets

---

## 🔎 7. Monitoring

- [ ] Deploy `postgres_exporter`
- [ ] Install Prometheus + Grafana (e.g., kube-prometheus-stack)
- [ ] Import dashboards: replication, slow queries, WAL stats
- [ ] Set alerts: WAL disk full, replication lag, etc.
- [ ] Tail operator logs for sync/failover behavior

---

## 🔁 8. High Availability

- [ ] Use `numberOfInstances: 3`
- [ ] Enable podAntiAffinity or topologySpreadConstraints
- [ ] Spread across 2–3 nodes
- [ ] Test failover and auto-promotion

---

## 🔄 9. Postgres Upgrades

- [ ] Use version like `"15"` or `"16"`
- [ ] Change version in CRD for rolling upgrade
- [ ] Test major upgrades in staging/dev
- [ ] Use `VACUUM FULL` post-upgrade if needed

---

## 📊 10. Connection Pooling

- [ ] Enable `enableConnectionPooler: true` in CRD
- [ ] Use `pgbouncer-rr` mode
- [ ] Use service name like `myapp-db-pooler.namespace.svc.cluster.local`
- [ ] Tune pooler PodDisruptionBudget

---

## 📘 Bonus: GitOps & Infra-as-Code

- [ ] Store all manifests in Git
- [ ] Use Kustomize or Helm for overlays
- [ ] Apply with ArgoCD, FluxCD or CI
- [ ] Version schema migrations (e.g., Goose, Liquibase)

---

## ✅ Summary

| Feature                  | Ready |
|--------------------------|-------|
| ARM Compatibility        | ✅    |
| Persistent Storage (CSI) | ✅    |
| WAL-G Backups            | ✅    |
| TLS + Secrets            | ✅    |
| HA Replication           | ✅    |
| Monitoring + Dashboards  | ✅    |
| Pooling + Performance    | ✅    |
| GitOps Ready             | ✅    |
