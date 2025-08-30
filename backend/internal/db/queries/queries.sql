-- name: UpsertUser :one
INSERT INTO users (sub, auth_provider_id, auth_provider, email, name, locale, profile_picture, verified_email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (email)
DO UPDATE SET
    name = COALESCE(EXCLUDED.name, users.name),
    profile_picture = COALESCE(EXCLUDED.profile_picture, users.profile_picture),
    locale = COALESCE(EXCLUDED.locale, users.locale),
    verified_email = COALESCE(EXCLUDED.verified_email, users.verified_email),
    last_login = NOW()
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetUserBySub :one
SELECT * FROM users
WHERE sub = $1;

-- name: GetUserBySubWithShops :one
SELECT 
    users.*,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'shop_id', shops.shop_id,
                'title', shops.title,
                'domain', shops.subdomain,
                'subdomain', shops.subdomain,
                'status', shops.status,
                'created_at', shops.created_at,
                'updated_at', shops.updated_at
            )
        ) FILTER (WHERE shops.shop_id IS NOT NULL), '[]'::jsonb
    )::jsonb AS shops
FROM users
LEFT JOIN shops ON users.user_id = shops.owner_id
WHERE users.sub = $1
GROUP BY users.user_id;

-- name: UpsertCustomer :one
INSERT INTO shop_customers (email, name, locale, profile_picture, verified_email, auth_provider, auth_provider_id, shop_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (email, shop_id)
DO UPDATE SET
    name = COALESCE(EXCLUDED.name, shop_customers.name),
    locale = COALESCE(EXCLUDED.locale, shop_customers.locale),
    profile_picture = COALESCE(EXCLUDED.profile_picture, shop_customers.profile_picture),
    verified_email = COALESCE(EXCLUDED.verified_email, shop_customers.verified_email),
    auth_provider = COALESCE(EXCLUDED.auth_provider, shop_customers.auth_provider),
    auth_provider_id = COALESCE(EXCLUDED.auth_provider_id, shop_customers.auth_provider_id),
    last_login = NOW()
RETURNING *;

-- name: GetCustomerByEmail :one
SELECT * FROM shop_customers
WHERE shop_customers.email = $1 AND shop_id = (SELECT shop_id FROM shops WHERE subdomain = $2);

-- name: GetCustomers :many
SELECT * FROM shop_customers
WHERE shop_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCustomersCount :one
SELECT COUNT(*) FROM shop_customers
WHERE shop_id = $1;

-- name: SearchCustomers :many
SELECT * FROM shop_customers
WHERE shop_id = $1 
AND (
    LOWER(name) LIKE LOWER($2) OR 
    LOWER(email) LIKE LOWER($2)
)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetCustomerById :one
SELECT * FROM shop_customers
WHERE shop_customer_id = $1 AND shop_id = $2;

-- name: UpdateCustomer :one
UPDATE shop_customers 
SET 
    name = COALESCE($3, name),
    locale = COALESCE($4, locale),
    profile_picture = COALESCE($5, profile_picture),
    verified_email = COALESCE($6, verified_email),
    last_login = NOW()
WHERE shop_customer_id = $1 AND shop_id = $2
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM shop_customers
WHERE shop_customer_id = $1 AND shop_id = $2;

-- name: GetCustomerOrders :many
SELECT o.*, oi.order_item_id, oi.product_variation_id, oi.quantity, oi.price as item_price
FROM orders o
LEFT JOIN order_items oi ON o.order_id = oi.order_id
WHERE o.customer_email = $1 AND o.shop_id = $2
ORDER BY o.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetLowStockVariants :many
SELECT pv.*, p.title as product_title
FROM product_variations pv
JOIN products p ON pv.product_id = p.product_id
WHERE pv.shop_id = $1 
AND pv.available_quantity <= $2
ORDER BY pv.available_quantity ASC;

-- name: GetProductVariation :one
SELECT * FROM product_variations
WHERE product_variation_id = $1 AND shop_id = $2;

-- name: UpdateVariantStock :one
UPDATE product_variations
SET available_quantity = $3,
    updated_at = NOW()
WHERE product_variation_id = $1 AND shop_id = $2
RETURNING *;

-- name: DeductVariantStock :one
UPDATE product_variations
SET available_quantity = GREATEST(available_quantity - $3, 0),
    updated_at = NOW()
WHERE product_variation_id = $1 AND shop_id = $2
RETURNING *;

-- name: AddVariantStock :one
UPDATE product_variations
SET available_quantity = available_quantity + $3,
    updated_at = NOW()
WHERE product_variation_id = $1 AND shop_id = $2
RETURNING *;

-- name: GetInventoryReport :many
SELECT 
    pv.product_variation_id,
    pv.product_id,
    p.title as product_title,
    pv.description as variant_description,
    pv.sku,
    pv.available_quantity,
    0 as reserved_quantity,
    pv.available_quantity as available_stock,
    pv.price,
    (pv.available_quantity * pv.price)::numeric(12,2) as stock_value,
    pv.updated_at,
    CASE 
        WHEN pv.available_quantity = 0 THEN 'OUT_OF_STOCK'
        WHEN pv.available_quantity <= $2 THEN 'LOW_STOCK'
        ELSE 'IN_STOCK'
    END as stock_status
FROM product_variations pv
JOIN products p ON pv.product_id = p.product_id
WHERE pv.shop_id = $1
ORDER BY pv.available_quantity ASC;

-- name: GetStockMovements :many
SELECT 
    sm.*,
    p.title as product_title,
    pv.description as variant_title,
    pv.sku
FROM stock_movements sm
JOIN product_variations pv ON sm.product_variation_id = pv.product_variation_id
JOIN products p ON pv.product_id = p.product_id
WHERE sm.shop_id = $1
ORDER BY sm.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateStockMovement :one
INSERT INTO stock_movements (
    product_variation_id, 
    shop_id, 
    movement_type, 
    quantity_change, 
    quantity_before, 
    quantity_after, 
    reference_id, 
    notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- Payment Status Management
-- name: UpdateOrderPaymentStatus :one
UPDATE orders
SET 
    payment_status = $3,
    transaction_id = $4,
    updated_at = NOW()
WHERE order_id = $1 AND shop_id = $2
RETURNING *;

-- name: GetOrderByTransactionID :one
SELECT * FROM orders
WHERE transaction_id = $1 AND shop_id = $2;

-- name: UpdateOrderStatusByTransactionID :one
UPDATE orders
SET 
    status = $3,
    payment_status = $4,
    updated_at = NOW()
WHERE transaction_id = $1 AND shop_id = $2
RETURNING *;
