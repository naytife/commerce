-- name: CreateCategory :one
INSERT INTO categories (slug, title, description, parent_id, shop_id, category_attributes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCategory :one
SELECT category_id, slug, title, description, created_at, updated_at, parent_id, category_attributes
FROM categories
WHERE shop_id = $1 AND category_id = $2;

-- name: UpdateCategory :one
UPDATE categories
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    parent_id = COALESCE(sqlc.narg('parent_id'), parent_id)
WHERE category_id = sqlc.arg('category_id')
RETURNING *;

-- name: GetCategories :many
SELECT category_id, slug, title, description, created_at, updated_at
FROM categories
WHERE shop_id = sqlc.arg('shop_id') AND category_id > sqlc.arg('after')
LIMIT sqlc.arg('limit');

-- name: GetCategoryChildren :many
SELECT category_id, slug, title, description, created_at, updated_at
FROM categories
WHERE shop_id = sqlc.arg('shop_id') AND parent_id = sqlc.arg('parent_id');

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

-- name: GetCategoryAttributes :one
SELECT category_attributes
FROM categories
WHERE category_id = sqlc.arg('category_id');