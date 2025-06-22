-- name: CreateDeployment :one
INSERT INTO shop_deployments (
    shop_id,
    template_name,
    template_version,
    status,
    deployment_type,
    message,
    started_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateDeploymentStatus :exec
UPDATE shop_deployments 
SET status = $2, message = $3, updated_at = CURRENT_TIMESTAMP
WHERE deployment_id = $1;

-- name: CompleteDeployment :exec
UPDATE shop_deployments 
SET status = $2, message = $3, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE deployment_id = $1;

-- name: GetDeploymentByID :one
SELECT * FROM shop_deployments 
WHERE deployment_id = $1;

-- name: GetLatestDeploymentByShop :one
SELECT * FROM shop_deployments 
WHERE shop_id = $1 
ORDER BY started_at DESC 
LIMIT 1;

-- name: GetDeploymentsByShop :many
SELECT * FROM shop_deployments 
WHERE shop_id = $1 
ORDER BY started_at DESC 
LIMIT $2 OFFSET $3;

-- name: GetShopCurrentTemplate :one
SELECT template_name, template_version, status, completed_at
FROM shop_deployments 
WHERE shop_id = $1 AND status = 'deployed'
ORDER BY completed_at DESC 
LIMIT 1;

-- name: IsShopDeployed :one
SELECT EXISTS(
    SELECT 1 FROM shop_deployments 
    WHERE shop_id = $1 AND status = 'deployed'
) as is_deployed;

-- name: UpdateShopLastDeployment :exec
UPDATE shops 
SET last_deployment_id = $2, updated_at = CURRENT_TIMESTAMP
WHERE shop_id = $1;

-- Data update tracking queries

-- name: CreateDataUpdate :one
INSERT INTO shop_data_updates (
    shop_id,
    data_type,
    status,
    changes_summary,
    started_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateDataUpdateStatus :exec
UPDATE shop_data_updates 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE update_id = $1;

-- name: CompleteDataUpdate :exec
UPDATE shop_data_updates 
SET status = $2, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE update_id = $1;

-- name: GetDataUpdateByID :one
SELECT * FROM shop_data_updates 
WHERE update_id = $1;

-- name: GetLatestDataUpdateByShop :one
SELECT * FROM shop_data_updates 
WHERE shop_id = $1 
ORDER BY started_at DESC 
LIMIT 1;

-- name: UpdateShopLastDataUpdate :exec
UPDATE shops 
SET last_data_update_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE shop_id = $1;

-- Deployment URL tracking

-- name: CreateDeploymentURL :one
INSERT INTO shop_deployment_urls (
    deployment_id,
    url_type,
    url
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetDeploymentURLs :many
SELECT * FROM shop_deployment_urls 
WHERE deployment_id = $1
ORDER BY url_type;
