-- name: CreateShop :one
INSERT INTO shops (owner_id, title, subdomain, email, currency_code, about, status, address, phone_number, seo_description, seo_keywords, seo_title, current_template)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE shop_id = $1;

-- name: DeleteShop :exec
DELETE FROM shops
WHERE shop_id = $1;

-- name: GetShopImages :one
SELECT * FROM shop_images
WHERE shop_id = $1;

-- name: CreateShopImages :one
INSERT INTO shop_images (
    favicon_url, 
    logo_url, 
    logo_url_dark, 
    banner_url, 
    banner_url_dark, 
    cover_image_url, 
    cover_image_url_dark,
    shop_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateShopImages :one
UPDATE shop_images
SET 
    favicon_url = COALESCE(sqlc.narg('favicon_url'), favicon_url),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    logo_url_dark = COALESCE(sqlc.narg('logo_url_dark'), logo_url_dark),
    banner_url = COALESCE(sqlc.narg('banner_url'), banner_url),
    banner_url_dark = COALESCE(sqlc.narg('banner_url_dark'), banner_url_dark),
    cover_image_url = COALESCE(sqlc.narg('cover_image_url'), cover_image_url),
    cover_image_url_dark = COALESCE(sqlc.narg('cover_image_url_dark'), cover_image_url_dark)
WHERE shop_id = sqlc.arg('shop_id')
RETURNING *;

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
    whatsapp_link = COALESCE(sqlc.narg('whatsapp_link'), whatsapp_link),
    whatsapp_phone_number = COALESCE(sqlc.narg('whatsapp_phone_number'), whatsapp_phone_number),
    facebook_link = COALESCE(sqlc.narg('facebook_link'), facebook_link),
    instagram_link = COALESCE(sqlc.narg('instagram_link'), instagram_link),
    seo_description = COALESCE(sqlc.narg('seo_description'), seo_description),
    seo_keywords = COALESCE(sqlc.narg('seo_keywords'), seo_keywords),
    seo_title = COALESCE(sqlc.narg('seo_title'), seo_title),
    address = COALESCE(sqlc.narg('address'), address),
    email = COALESCE(sqlc.narg('email'), email),
    current_template = COALESCE(sqlc.narg('current_template'), current_template),
    updated_at = NOW()
WHERE shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: GetShopBySubDomain :one
SELECT * FROM shops
WHERE subdomain = $1;

-- name: GetShopIDBySubDomain :one
SELECT shop_id FROM shops
WHERE subdomain = $1;