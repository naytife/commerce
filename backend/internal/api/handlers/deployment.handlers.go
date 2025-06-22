package handlers

// NOTE: This file previously contained deprecated deployment handlers and duplicate type definitions.
//
// DEPLOYMENT FUNCTIONALITY:
// - All deployment functionality has been migrated to proxy handlers in proxy.handlers.go
// - Use ProxyDeployStore, ProxyRedeployStore, ProxyDeploymentStatus, and ProxyUpdateStoreData
//
// TYPE DEFINITIONS:
// - All deployment-related types are now defined in internal/api/models/template.go
// - This ensures single source of truth and eliminates duplication
//
// MIGRATION COMPLETE:
// - Zero handlers remain in this file (successfully migrated to proxy pattern)
// - Zero types remain in this file (consolidated to models package)
// - File maintained for documentation and future reference
