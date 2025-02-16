
-- name: CreateProduct :one
INSERT INTO products ( title, description, status, product_type_id, shop_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

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

-- name: UpdateProduct :exec
UPDATE products
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE product_id = sqlc.arg('product_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = $1 AND shop_id = $2
RETURNING *;

-- name: UpsertProductVariants :batchexec
INSERT INTO product_variations (slug, description, price,sku, available_quantity, seo_description, seo_keywords, seo_title, product_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (slug, shop_id)
DO UPDATE SET
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    sku = EXCLUDED.sku,
    available_quantity = EXCLUDED.available_quantity,
    seo_description = EXCLUDED.seo_description,
    seo_keywords = EXCLUDED.seo_keywords,
    seo_title = EXCLUDED.seo_title
RETURNING *;

-- name: DeleteProductVariants :batchexec
DELETE FROM product_variations
WHERE shop_id = $1 AND product_id = $2
AND product_variation_id != ALL(sqlc.arg('product_variation_ids')::bigint[]);

-- name: GetProductVariants :many
SELECT * FROM product_variations
WHERE shop_id = $1 AND product_id = $2
ORDER BY product_variation_id;

-- name: GetProduct :one
SELECT 
    p.product_id,
    p.title,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Aggregate attributes separately to prevent duplication
    (
        SELECT COALESCE(
            json_agg(
                json_build_object(
                    'attribute_id', pa.attribute_id,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::json
        )
        FROM product_attribute_values pa
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    ) AS attributes,

    -- Aggregate variants separately to prevent duplication
    (
        SELECT COALESCE(
            json_agg(DISTINCT jsonb_build_object(
                'variation_id', pv.product_variation_id,
                'slug', pv.slug,
                'description', pv.description,
                'price', pv.price,
                'sku', pv.sku,
                'available_quantity', pv.available_quantity
            )) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::json
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    ) AS variants

FROM products p
WHERE p.product_id = $1 AND p.shop_id = $2;

-- name: GetProducts :many
SELECT 
    p.product_id,
    p.title,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Aggregate attributes separately to prevent duplication
    (
        SELECT COALESCE(
            json_agg(
                json_build_object(
                    'attribute_id', pa.attribute_id,
                    'attribute_title', a.title,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::json
        )
        FROM product_attribute_values pa
        LEFT JOIN attributes a ON a.attribute_id = pa.attribute_id
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    ) AS attributes,

    -- Aggregate variants separately to prevent duplication
    (
        SELECT COALESCE(
            json_agg(DISTINCT jsonb_build_object(
                'variation_id', pv.product_variation_id,
                'slug', pv.slug,
                'description', pv.description,
                'price', pv.price,
                'sku', pv.sku,
                'available_quantity', pv.available_quantity
            )) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::json
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    ) AS variants

FROM products p
WHERE p.shop_id = sqlc.arg('shop_id') 
AND p.product_id > sqlc.arg('after')
ORDER BY p.product_id
LIMIT sqlc.arg('limit');

-- name: GetProductsByType :many
SELECT 
    p.product_id,
    p.title,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Aggregate attributes separately
    (
        SELECT COALESCE(
            json_agg(
                json_build_object(
                    'attribute_id', pa.attribute_id,
                    'attribute_title', a.title,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::json
        )
        FROM product_attribute_values pa
        LEFT JOIN attributes a ON a.attribute_id = pa.attribute_id
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    ) AS attributes,

    -- Aggregate variants separately
    (
        SELECT COALESCE(
            json_agg(DISTINCT jsonb_build_object(
                'variation_id', pv.product_variation_id,
                'slug', pv.slug,
                'description', pv.description,
                'price', pv.price,
                'sku', pv.sku,
                'available_quantity', pv.available_quantity
            )) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::json
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    ) AS variants

FROM products p
WHERE p.shop_id = sqlc.arg('shop_id') 
AND p.product_type_id = sqlc.arg('product_type_id')
AND p.product_id > sqlc.arg('after')
ORDER BY p.product_id
LIMIT sqlc.arg('limit');