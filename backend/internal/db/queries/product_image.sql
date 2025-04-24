-- name: CreateProductImage :one
INSERT INTO product_images (url, alt, product_id, shop_id) 
VALUES ($1, $2, $3, $4)
RETURNING product_image_id, url, alt, product_id, shop_id;

-- name: GetProductImages :many
SELECT product_image_id, url, alt, product_id, shop_id
FROM product_images
WHERE product_id = $1 AND shop_id = $2
ORDER BY product_image_id;

-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE product_image_id = $1 AND shop_id = $2;

-- name: DeleteAllProductImages :exec
DELETE FROM product_images
WHERE product_id = $1 AND shop_id = $2; 