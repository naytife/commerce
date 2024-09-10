-- name: UpsertUser :one
INSERT INTO users (auth0_sub, email, name, profile_picture_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (auth0_sub)
DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name, profile_picture_url = EXCLUDED.profile_picture_url
RETURNING user_id, auth0_sub, email, name, profile_picture_url;

-- name: GetUser :one
SELECT * FROM users
WHERE auth0_sub = $1;

-- name: CreateShop :one
INSERT INTO shops (owner_id, title, default_domain, favicon_url, currency_code, about, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE shop_id = $1;

-- name: GetShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;

-- name: UpdateShop :one
UPDATE shops
SET title = $2, favicon_url = $3, currency_code = $4, about = $5, status = $6
WHERE default_domain = $1
RETURNING *;

-- name: DeleteShop :exec
DELETE FROM shops
WHERE shop_id = $1;

-- name: GetShopByDomain :one
SELECT * FROM shops
WHERE default_domain = $1;

-- name: GetWhatsappsByShop :many
SELECT * FROM whatsapps
WHERE shop_id = $1;
