
-- name: CreateShop :one
INSERT INTO shops (owner_id, title, domain,email, currency_code, about, status, address,phone_number, seo_description, seo_keywords, seo_title)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE shop_id = $1;

-- name: GetShopImages :one
SELECT * FROM shop_images
WHERE shop_id = $1;

-- name: GetShopWhatsApp :one
SELECT * FROM whatsapps
WHERE shop_id = $1;

-- name: GetShopFacebook :one
SELECT * FROM facebooks
WHERE shop_id = $1;

-- name: GetShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;

-- name: UpdateShop :one
UPDATE shops
SET 
    title = COALESCE(sqlc.narg('title'), title),
    currency_code = COALESCE(sqlc.narg('currency_code'), currency_code),
    about = COALESCE(sqlc.narg('about'), about),
    status = COALESCE(sqlc.narg('status'), status),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    seo_description = COALESCE(sqlc.narg('seo_description'), seo_description),
    seo_keywords = COALESCE(sqlc.narg('seo_keywords'), seo_keywords),
    seo_title = COALESCE(sqlc.narg('seo_title'), seo_title),
    address = COALESCE(sqlc.narg('address'), address),
    email = COALESCE(sqlc.narg('email'), email)
WHERE shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: UpsertShopWhatsapp :one
INSERT INTO whatsapps (shop_id, url, phone_number)
VALUES ($1, $2, $3)
ON CONFLICT (shop_id)
DO UPDATE SET url = EXCLUDED.url, phone_number = EXCLUDED.phone_number
RETURNING *;

-- name: UpsertShopFacebook :one
INSERT INTO facebooks (shop_id, url, handle)
VALUES ($1, $2, $3)
ON CONFLICT (shop_id)
DO UPDATE SET url = EXCLUDED.url, handle = EXCLUDED.handle
RETURNING *;

-- name: GetShopByDomain :one
SELECT * FROM shops
WHERE domain = $1;

-- name: GetShopIDByDomain :one
SELECT shop_id FROM shops
WHERE domain = $1;