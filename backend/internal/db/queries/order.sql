-- name: CreateOrder :one
INSERT INTO orders (
    status, amount, discount, shipping_cost, tax,
    shipping_address, payment_method, payment_status,
    shipping_method, shipping_status, transaction_id,
    username, shop_customer_id, shop_id, customer_name,
    customer_email, customer_phone
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8,
    $9, $10, $11,
    $12, $13, $14, $15,
    $16, $17
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE order_id = $1 AND shop_id = $2;

-- name: ListOrders :many
SELECT * FROM orders
WHERE shop_id = $1
ORDER BY order_id
LIMIT $2 OFFSET $3;

-- name: UpdateOrder :exec
UPDATE orders
SET
    status = $1,
    amount = $2,
    discount = $3,
    shipping_cost = $4,
    tax = $5,
    shipping_address = $6,
    payment_method = $7,
    payment_status = $8,
    shipping_method = $9,
    shipping_status = $10,
    transaction_id = $11,
    username = $12,
    customer_name = $13,
    customer_email = $14,
    customer_phone = $15,
    updated_at = NOW()
WHERE order_id = $16 AND shop_id = $17;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1 AND shop_id = $2;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    quantity, price, product_variation_id, order_id, shop_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOrderItemsByOrder :many
SELECT * FROM order_items
WHERE order_id = $1 AND shop_id = $2
ORDER BY order_item_id;

-- name: UpdateOrderItem :exec
UPDATE order_items
SET
    quantity = $1,
    price = $2,
    updated_at = NOW()
WHERE order_item_id = $3 AND shop_id = $4
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE order_item_id = $1 AND shop_id = $2;

-- name: DeleteOrderItemsByOrder :exec
DELETE FROM order_items
WHERE order_id = $1 AND shop_id = $2; 