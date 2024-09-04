-- name: CreateShop :one
INSERT INTO shops (owner_id, title, default_domain, favicon_url, currency_code, about, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE id = $1;

-- name: GetShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;