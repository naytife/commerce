
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

-- name: UpdateProduct :one
UPDATE products
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description)
WHERE product_id = sqlc.arg('product_id')
RETURNING *;

-- name: CreateProductAllowedAttribute :one
UPDATE products
SET allowed_attributes = jsonb_set(
    COALESCE(allowed_attributes, '{}'), 
    ARRAY[UPPER(sqlc.arg('title'))::text], 
    to_jsonb(sqlc.arg('data_type')::text)
)
WHERE product_id = sqlc.arg('product_id')
RETURNING allowed_attributes;

-- name: DeleteProductAllowedAttribute :one
UPDATE products
SET allowed_attributes = allowed_attributes - UPPER(sqlc.arg('attribute')::text)
WHERE product_id = sqlc.arg('product_id')
RETURNING allowed_attributes;

-- name: UpsertProductVariation :batchmany
INSERT INTO product_variations (slug, description, price, available_quantity, attributes, seo_description, seo_keywords, seo_title, product_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (slug, shop_id)
DO UPDATE SET
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    available_quantity = EXCLUDED.available_quantity,
    attributes = EXCLUDED.attributes,
    seo_description = EXCLUDED.seo_description,
    seo_keywords = EXCLUDED.seo_keywords,
    seo_title = EXCLUDED.seo_title
RETURNING *;


-- name: DeleteProductVariations :batchexec
DELETE FROM product_variations
WHERE shop_id = $1 AND product_id = $2
AND product_variation_id != ALL(sqlc.arg('product_variation_ids')::bigint[]);

-- name: GetProductVariations :many
SELECT * FROM product_variations
WHERE shop_id = $1 AND product_id = $2
ORDER BY product_variation_id;