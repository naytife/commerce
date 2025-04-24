-- name: CreateProduct :one
INSERT INTO products ( slug, title, description, status, product_type_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6)
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

-- name: UpsertProductVariants :batchmany
INSERT INTO product_variations (
    description, price, sku, available_quantity,
    seo_description, seo_keywords, seo_title,
    product_id, shop_id, is_default
)
VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8,
    $9, $10
)
ON CONFLICT (sku, shop_id)
DO UPDATE SET
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    available_quantity = EXCLUDED.available_quantity,
    seo_description = EXCLUDED.seo_description,
    seo_keywords = EXCLUDED.seo_keywords,
    seo_title = EXCLUDED.seo_title
RETURNING *;

-- name: DeleteProductVariants :exec
DELETE FROM product_variations
WHERE shop_id = sqlc.arg('shop_id') AND product_id = sqlc.arg('product_id')
AND product_variation_id != ALL(sqlc.arg('product_variation_ids')::bigint[]);

-- name: GetProductVariants :many
SELECT * FROM product_variations
WHERE shop_id = $1 AND product_id = $2
ORDER BY product_variation_id;

-- name: GetProduct :one
SELECT 
    p.product_id,
    p.title,
    p.slug,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Aggregate attributes separately to prevent duplication
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'attribute_id', pa.attribute_id,
                    'title', a.title,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_attribute_values pa
        LEFT JOIN attributes a ON a.attribute_id = pa.attribute_id
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    )::jsonb AS attributes,

    -- Aggregate variants separately to prevent duplication
      (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'variation_id', pv.product_variation_id,
                    'description', pv.description,
                    'price', pv.price,
                    'sku', pv.sku,
                    'available_quantity', pv.available_quantity,
                    'is_default', pv.is_default,
                    'attributes', (
                        SELECT COALESCE(
                            jsonb_agg(
                                jsonb_build_object(
                                    'attribute_id', pva.attribute_id,
                                    'title', a.title,
                                    'attribute_option_id', pva.attribute_option_id,
                                    'value', COALESCE(ao.value, pva.value)
                                )
                            ) FILTER (WHERE pva.attribute_id IS NOT NULL),
                            '[]'::jsonb
                        )
                        FROM product_variation_attribute_values pva
                        LEFT JOIN attributes a ON a.attribute_id = pva.attribute_id
                        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pva.attribute_option_id
                        WHERE pva.product_variation_id = pv.product_variation_id
                    )
                )
            ) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    )::jsonb AS variants,
    
    -- Aggregate product images separately
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', pi.product_image_id,
                    'url', pi.url,
                    'alt', pi.alt
                )
            ) FILTER (WHERE pi.product_image_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_images pi
        WHERE pi.product_id = p.product_id
    )::jsonb AS images

FROM products p
WHERE p.product_id = $1 AND p.shop_id = $2;

-- name: GetProducts :many
SELECT 
    p.product_id,
    p.slug,
    p.title,
    p.description,
    p.status,
    p.category_id,
    p.updated_at,
    p.created_at,

    -- Product attributes
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'attribute_id', pa.attribute_id,
                    'title', a.title,
                    'attribute_option_id', pa.attribute_option_id,
                    'value', COALESCE(ao.value, pa.value)
                )
            ) FILTER (WHERE pa.attribute_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_attribute_values pa
        LEFT JOIN attributes a ON a.attribute_id = pa.attribute_id
        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pa.attribute_option_id
        WHERE pa.product_id = p.product_id
    )::jsonb AS attributes,

    -- Product variants with embedded attributes
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'variation_id', pv.product_variation_id,
                    'description', pv.description,
                    'price', pv.price,
                    'sku', pv.sku,
                    'available_quantity', pv.available_quantity,
                    'is_default', pv.is_default,
                    'attributes', (
                        SELECT COALESCE(
                            jsonb_agg(
                                jsonb_build_object(
                                    'attribute_id', pva.attribute_id,
                                    'title', a.title,
                                    'attribute_option_id', pva.attribute_option_id,
                                    'value', COALESCE(ao.value, pva.value)
                                )
                            ) FILTER (WHERE pva.attribute_id IS NOT NULL),
                            '[]'::jsonb
                        )
                        FROM product_variation_attribute_values pva
                        LEFT JOIN attributes a ON a.attribute_id = pva.attribute_id
                        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pva.attribute_option_id
                        WHERE pva.product_variation_id = pv.product_variation_id
                    )
                )
            ) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    )::jsonb AS variants,
    
    -- Product images
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', pi.product_image_id,
                    'url', pi.url,
                    'alt', pi.alt
                )
            ) FILTER (WHERE pi.product_image_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_images pi
        WHERE pi.product_id = p.product_id
    )::jsonb AS images

FROM products p
WHERE p.shop_id = sqlc.arg('shop_id') 
AND p.product_id > sqlc.arg('after')
ORDER BY p.product_id
LIMIT sqlc.arg('limit');

-- name: GetProductsByType :many
SELECT 
    p.product_id,
    p.slug,
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
                    'title', a.title,
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
            jsonb_agg(
                jsonb_build_object(
                    'variation_id', pv.product_variation_id,
                    'description', pv.description,
                    'price', pv.price,
                    'sku', pv.sku,
                    'available_quantity', pv.available_quantity,
                    'is_default', pv.is_default,
                    'attributes', (
                        SELECT COALESCE(
                            jsonb_agg(
                                jsonb_build_object(
                                    'attribute_id', pva.attribute_id,
                                    'title', a.title,
                                    'attribute_option_id', pva.attribute_option_id,
                                    'value', COALESCE(ao.value, pva.value)
                                )
                            ) FILTER (WHERE pva.attribute_id IS NOT NULL),
                            '[]'::jsonb
                        )
                        FROM product_variation_attribute_values pva
                        LEFT JOIN attributes a ON a.attribute_id = pva.attribute_id
                        LEFT JOIN attribute_options ao ON ao.attribute_option_id = pva.attribute_option_id
                        WHERE pva.product_variation_id = pv.product_variation_id
                    )
                )
            ) FILTER (WHERE pv.product_variation_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_variations pv
        WHERE pv.product_id = p.product_id
    )::jsonb AS variants,

    -- Aggregate product images
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', pi.product_image_id,
                    'url', pi.url,
                    'alt', pi.alt
                )
            ) FILTER (WHERE pi.product_image_id IS NOT NULL),
            '[]'::jsonb
        )
        FROM product_images pi
        WHERE pi.product_id = p.product_id
    )::jsonb AS images

FROM products p
WHERE p.shop_id = sqlc.arg('shop_id') 
AND p.product_type_id = sqlc.arg('product_type_id')
AND p.product_id > sqlc.arg('after')
ORDER BY p.product_id
LIMIT sqlc.arg('limit');

-- name: UpdateProductVariationSku :one
UPDATE product_variations
SET sku = $2
WHERE product_variation_id = $1 AND shop_id = $3
RETURNING *;

-- name: CreateProductVariation :one
INSERT INTO product_variations (
    description, price, sku, available_quantity,
    seo_description, seo_keywords, seo_title,
    product_id, shop_id, is_default
)
VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8,
    $9, $10
)
RETURNING *;

-- name: UpdateProductVariation :one
UPDATE product_variations
SET 
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    available_quantity = COALESCE(sqlc.narg('available_quantity'), available_quantity),
    seo_description = COALESCE(sqlc.narg('seo_description'), seo_description),
    seo_keywords = COALESCE(sqlc.narg('seo_keywords'), seo_keywords),
    seo_title = COALESCE(sqlc.narg('seo_title'), seo_title),
    is_default = COALESCE(sqlc.narg('is_default'), is_default),
    updated_at = NOW()
WHERE product_variation_id = sqlc.arg('product_variation_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: GetProductById :one
SELECT * FROM products
WHERE product_id = $1 AND shop_id = $2;

-- name: GetProductTypeByProduct :one
SELECT pt.*
FROM products p
JOIN product_types pt ON p.product_type_id = pt.product_type_id
WHERE p.product_id = sqlc.arg('product_id') AND p.shop_id = sqlc.arg('shop_id');