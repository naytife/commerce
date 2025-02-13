-- name: CreateProductType :one
INSERT INTO product_types ( title, shippable, digital, shop_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProductTypes :many
SELECT * FROM product_types WHERE shop_id = $1;

-- name: GetProductType :one
SELECT * FROM product_types WHERE product_type_id = $1 AND shop_id = $2;

-- name: UpdateProductType :one
UPDATE product_types
SET 
    title = COALESCE(sqlc.narg('title'), title),
    shippable = COALESCE(sqlc.narg('shippable'), shippable),
    digital = COALESCE(sqlc.narg('digital'), digital)
WHERE product_type_id = sqlc.arg('product_type_id') AND shop_id = sqlc.arg('shop_id')
RETURNING *;

-- name: DeleteProductType :one
DELETE FROM product_types
WHERE product_type_id = $1 AND shop_id = $2
RETURNING *;

