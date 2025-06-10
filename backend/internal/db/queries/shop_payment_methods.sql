-- name: GetShopPaymentMethods :many
SELECT * FROM shop_payment_methods
WHERE shop_id = $1
ORDER BY method_type;

-- name: GetShopPaymentMethod :one
SELECT * FROM shop_payment_methods
WHERE shop_id = $1 AND method_type = $2;

-- name: UpsertShopPaymentMethod :one
INSERT INTO shop_payment_methods (shop_id, method_type, is_enabled, attributes)
VALUES ($1, $2, $3, $4)
ON CONFLICT (shop_id, method_type)
DO UPDATE SET
    is_enabled = EXCLUDED.is_enabled,
    attributes = EXCLUDED.attributes,
    updated_at = NOW()
RETURNING *;

-- name: UpdateShopPaymentMethodStatus :one
UPDATE shop_payment_methods
SET is_enabled = $3, updated_at = NOW()
WHERE shop_id = $1 AND method_type = $2
RETURNING *;

-- name: DeleteShopPaymentMethod :exec
DELETE FROM shop_payment_methods
WHERE shop_id = $1 AND method_type = $2;
