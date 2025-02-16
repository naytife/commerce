-- name: CreateAttribute :one
INSERT INTO attributes (title, data_type, unit, required, applies_to, product_type_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAttributes :many
SELECT * FROM attributes WHERE product_type_id = $1 AND shop_id = $2;

-- name: GetProductsAttributes :many
SELECT * FROM attributes WHERE applies_to = 'Product' AND product_type_id = $1 AND shop_id = $2;

-- name: GetVariationsAttributes :many
SELECT * FROM attributes WHERE applies_to = 'ProductVariation' AND product_type_id = $1 AND shop_id = $2;

-- name: GetAttribute :one
SELECT * FROM attributes WHERE attribute_id = $1 AND shop_id = $2;

-- name: UpdateAttribute :one
UPDATE attributes
SET 
    title = COALESCE(sqlc.narg('title'), title),
    data_type = COALESCE(sqlc.narg('data_type'), data_type),
    unit = COALESCE(sqlc.narg('unit'), unit),
    required = COALESCE(sqlc.narg('required'), required),
    applies_to = COALESCE(sqlc.narg('applies_to'), applies_to)
WHERE attribute_id = sqlc.arg('attribute_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: DeleteAttribute :one
DELETE FROM attributes
WHERE attribute_id = $1 AND shop_id = $2
RETURNING *;

-- name: CreateAttributeOption :one
INSERT INTO attribute_options (value, shop_id, attribute_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAttributeOptions :many
SELECT * FROM attribute_options WHERE attribute_id = $1 AND shop_id = $2;

-- name: GetAttributeOption :one
SELECT * FROM attribute_options WHERE attribute_option_id = $1 AND shop_id = $2;

-- name: UpdateAttributeOption :one
UPDATE attribute_options
SET 
    value = COALESCE(sqlc.narg('value'), value)
WHERE attribute_option_id = sqlc.arg('attribute_option_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: DeleteAttributeOption :one
DELETE FROM attribute_options
WHERE attribute_option_id = $1 AND shop_id = $2
RETURNING *;

-- name: BatchUpsertProductAttributeValues :batchexec
INSERT INTO product_attribute_values (value, attribute_option_id, product_id, attribute_id, shop_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (product_id, attribute_id, shop_id) 
DO UPDATE SET 
    value = EXCLUDED.value,
    attribute_option_id = EXCLUDED.attribute_option_id
WHERE product_attribute_values.value IS DISTINCT FROM EXCLUDED.value  
   OR product_attribute_values.attribute_option_id IS DISTINCT FROM EXCLUDED.attribute_option_id
RETURNING *;

-- name: BatchDeleteProductAttributeValues :batchexec
DELETE FROM product_attribute_values 
WHERE product_id = $1 
AND shop_id = $2
AND attribute_id NOT IN (SELECT UNNEST($3::int[]));

-- name: GetProductAttributeValues :many
SELECT 
    pav.product_id,
    pav.attribute_id,
    pav.shop_id,
    pav.attribute_option_id,
    COALESCE(ao.value, pav.value) as value
FROM product_attribute_values pav
LEFT JOIN attribute_options ao ON ao.attribute_option_id = pav.attribute_option_id
WHERE pav.product_id = $1 AND pav.shop_id = $2;
