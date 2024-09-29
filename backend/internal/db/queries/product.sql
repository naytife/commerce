
-- name: CreateProduct :one
INSERT INTO products ( title, description, category_id, shop_id, allowed_attributes, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProducts :many
SELECT product_id, title, description, created_at, updated_at, status, category_id
FROM products
WHERE shop_id = sqlc.arg('shop_id') AND product_id > sqlc.arg('after')
LIMIT sqlc.arg('limit');

-- name: GetProductsByCategory :many
SELECT product_id, title, description, created_at, updated_at, status, category_id
FROM products
WHERE category_id = sqlc.arg('category_id') AND product_id > sqlc.arg('after')
LIMIT sqlc.arg('limit');

-- name: GetProduct :one
SELECT product_id, title, description, created_at, updated_at, status, category_id
FROM products
WHERE shop_id = $1 AND product_id = $2;

-- name: GetProductAllowedAttributes :one
SELECT allowed_attributes
FROM products
WHERE product_id = sqlc.arg('product_id');