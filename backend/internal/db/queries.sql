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
INSERT INTO shops (owner_id, title, domain, favicon_url,logo_url,email, currency_code, about, status, address,phone_number, seo_description, seo_keywords, seo_title)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetShop :one
SELECT * FROM shops
WHERE shop_id = $1;

-- name: GetShopsByOwner :many
SELECT * FROM shops
WHERE owner_id = $1;

-- name: UpdateShop :one
UPDATE shops
SET 
    title = COALESCE(sqlc.narg('title'), title),
    favicon_url = COALESCE(sqlc.narg('favicon_url'), favicon_url),
    currency_code = COALESCE(sqlc.narg('currency_code'), currency_code),
    about = COALESCE(sqlc.narg('about'), about),
    status = COALESCE(sqlc.narg('status'), status),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    seo_description = COALESCE(sqlc.narg('seo_description'), seo_description),
    seo_keywords = COALESCE(sqlc.narg('seo_keywords'), seo_keywords),
    seo_title = COALESCE(sqlc.narg('seo_title'), seo_title),
    address = COALESCE(sqlc.narg('address'), address),
    email = COALESCE(sqlc.narg('email'), email)
WHERE domain = sqlc.arg('domain')
RETURNING *;

-- name: GetShopByDomain :one
SELECT * FROM shops
WHERE domain = $1;

-- name: GetShopIDByDomain :one
SELECT shop_id FROM shops
WHERE domain = $1;

-- name: CreateShopCategory :one
INSERT INTO categories (slug, title, description, parent_id, shop_id, category_attributes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetShopCategory :one
SELECT * FROM categories
WHERE category_id = $1;

-- name: UpdateShopCategory :one
UPDATE categories
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    parent_id = COALESCE(sqlc.narg('parent_id'), parent_id)
WHERE category_id = sqlc.arg('category_id')
RETURNING *;

-- name: GetShopCategories :many
SELECT * FROM categories;

-- name: CreateCategoryAttribute :one
UPDATE categories
SET category_attributes = jsonb_set(
    COALESCE(category_attributes, '{}'), 
    ARRAY[UPPER(sqlc.arg('title'))::text], 
    to_jsonb(sqlc.arg('data_type')::text)
)
WHERE category_id = sqlc.arg('category_id')
RETURNING category_attributes;

-- name: DeleteCategoryAttribute :one
UPDATE categories
SET category_attributes = category_attributes - sqlc.arg('attribute')::text
WHERE category_id = sqlc.arg('category_id')
RETURNING category_attributes;

-- name: CreateProduct :one
INSERT INTO products ( title, description, category_id, shop_id, allowed_attributes, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;