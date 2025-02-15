
-- name: CreateProduct :one
INSERT INTO products ( title, description, status, product_type_id, shop_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProducts :many
SELECT 
    p.product_id, 
    p.title, 
    p.description, 
    p.created_at, 
    p.updated_at, 
    p.status
FROM 
    products p
LEFT JOIN 
    categories c ON p.category_id = c.category_id
WHERE 
    p.shop_id = sqlc.arg('shop_id')
    AND p.product_id > sqlc.arg('after')
LIMIT sqlc.arg('limit');


-- SELECT 
--     p.product_id, 
--     p.title, 
--     p.description, 
--     p.created_at, 
--     p.updated_at, 
--     p.category_id
-- FROM 
--     products p
-- LEFT JOIN 
--     categories c ON p.category_id = c.category_id
-- WHERE 
--     p.shop_id = sqlc.arg('shop_id')
--     AND p.product_id > sqlc.arg('after')
-- LIMIT sqlc.arg('limit');

-- name: GetProductsByCategory :many
SELECT 
    p.product_id, 
    p.title, 
    p.description, 
    p.created_at, 
    p.updated_at, 
    p.category_id
FROM 
    products p
LEFT JOIN 
    categories c ON p.category_id = c.category_id
WHERE 
    p.category_id = sqlc.arg('category_id') 
    AND p.product_id > sqlc.arg('after')
LIMIT sqlc.arg('limit');

-- name: GetProduct :one
SELECT 
    p.product_id, 
    p.title, 
    p.description, 
    p.created_at, 
    p.updated_at, 
    p.category_id,
    p.status
FROM 
    products p
LEFT JOIN 
    categories c ON p.category_id = c.category_id
WHERE 
    p.shop_id = $1 
    AND p.product_id = $2;

-- xname: GetProductAllowedAttributes :one
-- SELECT 
--     (p.allowed_attributes || COALESCE(c.category_attributes, '{}'))::jsonb AS allowed_attributes
-- FROM 
--     products p
-- LEFT JOIN 
--     categories c ON p.category_id = c.category_id
-- WHERE 
--     p.product_id = sqlc.arg('product_id');

-- name: UpdateProduct :one
UPDATE products
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE product_id = sqlc.arg('product_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE product_id = $1 AND shop_id = $2
RETURNING *;

-- xname: CreateProductAllowedAttribute :one
-- UPDATE products
-- SET allowed_attributes = jsonb_set(
--     COALESCE(allowed_attributes, '{}'), 
--     ARRAY[UPPER(sqlc.arg('title'))::text], 
--     to_jsonb(sqlc.arg('data_type')::text)
-- )
-- WHERE product_id = sqlc.arg('product_id')
-- RETURNING allowed_attributes;

-- xname: DeleteProductAllowedAttribute :one
-- UPDATE products
-- SET allowed_attributes = allowed_attributes - UPPER(sqlc.arg('attribute')::text)
-- WHERE product_id = sqlc.arg('product_id')
-- RETURNING allowed_attributes;

-- name: UpsertProductVariation :batchmany
INSERT INTO product_variations (slug, description, price, available_quantity, seo_description, seo_keywords, seo_title, product_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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