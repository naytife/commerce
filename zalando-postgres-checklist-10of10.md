
# ✅ Zalando Postgres Operator 10/10 Checklist (k3s + ARM/x86 + Cloud)

This checklist guides you to deploy a **production-grade PostgreSQL cluster** using the **Zalando Postgres Operator** on your **k3s cluster running on ARM or x86 nodes**, hosted on **your current infrastructure** (no longer OCI). It prioritizes **reliability** without overengineering.

---

## 🔧 1. Cluster & Node Setup

- [x] Use ARM or x86 VMs (e.g., 4-core / 24GB RAM recommended)
- [x] Provision 2–3 nodes (for HA / PodAntiAffinity)
- [x] Use a reliable CSI plugin compatible with your cloud provider
- [x] Ensure nodes support required `linux/arm64` or `amd64` arch
- [x] Enable persistent storage with CSI driver
- [x] Install k3s with `--disable traefik` if using another ingress
- [x] Confirm PVC + CSIDriver support is active
- [x] Label nodes `postgres=true` if isolating DB workloads

---

## 🛠️ 2. Install Zalando Postgres Operator

- [x] Clone [Zalando Operator GitHub](https://github.com/zalando/postgres-operator)
- [x] Apply CRDs: `kubectl apply -f manifests/01-crds.yaml`
- [x] Deploy operator: `kubectl apply -f manifests/operator.yaml`
- [x] Create namespace (e.g., `database`)
- [x] Use multi-arch image like `ghcr.io/zalando/spilo-15:3.1-p1`
- [x] Confirm operator pod is running

---

## 🧠 3. Deploy PostgreSQL Cluster

- [x] Create `Postgresql` CRD with `numberOfInstances: 3`
- [x] Set `enableMasterLoadBalancer: false`
- [x] Define users, DBs, resources in the CRD
- [x] Use appropriate storage class
- [x] Use `teamId` to logically group DBs

---

## 💾 4. Persistent Storage

- [x] Use dynamic PVCs backed by CSI volumes
- [x] Ensure PVCs are bound and healthy
- [x] Validate data persists after pod/node restarts

---

## 🔐 5. Security & Secrets

- [x] Use Kubernetes Secrets for credentials
- [x] Enable TLS (internal + client-side)
- [x] Integrate with `cert-manager` for TLS automation if needed
- [x] Restrict RBAC access to Secrets

---

## 🔁 6. High Availability

- [x] Use `numberOfInstances: 3`
- [x] Enable `podAntiAffinity` or `topologySpreadConstraints`
- [x] Spread across 2–3 nodes
- [x] Test failover and auto-promotion

---

## 🔄 7. Postgres Upgrades

- [x] Use version like `"17"`
- [x] Change version in CRD for rolling upgrade
- [x] Test major upgrades in local/staging
- [x] Use `VACUUM FULL` post-upgrade if needed

---

## 📊 8. Connection Pooling

- [x] Enable `enableConnectionPooler: true` in CRD
- [x] Use `pgbouncer-rr` mode
- [x] Use service name like `myapp-db-pooler.namespace.svc.cluster.local`
- [x] Tune pooler `max_client_conn`, `default_pool_size`, etc.

---

## 📘 9. GitOps & Infra-as-Code

- [x] Store all manifests in Git
- [x] Use Kustomize or Helm for overlays
- [x] Apply with ArgoCD, FluxCD or CI
- [x] Version schema migrations (e.g., Goose, Liquibase)

---

## 🔁 10. Backups & PITR with WAL-G

- [x] Set up an S3-compatible object storage bucket (use R2)
- [x] Create Kubernetes `Secret` for WAL-G credentials (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `WALG_S3_PREFIX`)
- [x] Enable WAL archiving via `envFrom` in `Postgresql` CRD
- [x] Schedule logical backups with retention (e.g., cron: `"30 3 * * *"`)
- [x] Test a backup and full restore in dev

---

## 📊 11. Monitoring & Alerts

- [x] Deploy `postgres_exporter` (sidecar or DaemonSet)
- [x] Use Prometheus and Grafana (e.g., via kube-prometheus-stack)
- [x] Monitor key metrics: replication lag, WAL write speed, query duration
- [x] Configure alerts: replication lag > 5s, disk pressure, OOM
- [x] Integrate alerts to email, Slack, or PagerDuty

---

## 🔄 12. Maintenance & Upgrade Strategy

- [x] Test version upgrades on staging (`postgresql.version`)
- [x] Automate minor version upgrades via GitOps pipeline
- [x] Use `podDisruptionBudgets` to protect HA availability
- [x] Use anti-affinity and readiness/liveness probes
- [x] Document a disaster recovery process from backup

---

## ✅ Summary

| Feature                  | Ready |
|--------------------------|-------|
| ARM/x86 Compatibility    | ✅    |
| Persistent Storage (CSI) | ✅    |
| WAL-G Backups            | ✅    |
| TLS + Secrets            | ✅    |
| HA Replication           | ✅    |
| Monitoring + Dashboards  | ✅    |
| Pooling + Performance    | ✅    |
| GitOps Ready             | ✅    |
| Tested Recovery          | ✅    |
