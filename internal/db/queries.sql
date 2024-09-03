-- name: CreateShop :one
INSERT INTO shops (title, default_domain, favicon_url, currency_code, about)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE id = $1;

-- name: GetShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;