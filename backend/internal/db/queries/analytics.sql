-- name: GetSalesSummary :one
SELECT
  COALESCE(SUM(revenue), 0) AS total_sales,
  COALESCE(SUM(total_orders), 0) AS total_orders,
  COALESCE(SUM(revenue)/NULLIF(SUM(total_orders),0),0) AS average_order_value
FROM daily_sales
WHERE shop_id = $1 AND day BETWEEN $2 AND $3;

-- name: GetOrdersOverTime :many
SELECT
  DATE_TRUNC($1, day)::date AS period,
  COALESCE(SUM(total_orders), 0) AS order_count
FROM daily_sales
WHERE shop_id = $2 AND day BETWEEN $3 AND $4
GROUP BY period
ORDER BY period ASC;

-- name: GetTopProducts :many
SELECT
  oi.product_variation_id,
  p.title AS product_name,
  COALESCE(SUM(oi.quantity), 0) AS units_sold,
  COALESCE(SUM(oi.quantity * oi.price), 0) AS revenue
FROM order_items oi
JOIN orders o ON oi.order_id = o.order_id
JOIN product_variations pv ON oi.product_variation_id = pv.product_variation_id
JOIN products p ON pv.product_id = p.product_id
WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed'
GROUP BY oi.product_variation_id, p.title
ORDER BY units_sold DESC
LIMIT $4;

-- name: GetCustomerSummaryNewReturning :one
SELECT
  COUNT(DISTINCT CASE WHEN o.created_at BETWEEN $2 AND $3 THEN o.customer_email END) AS new_customers,
  COUNT(DISTINCT CASE WHEN o.created_at < $2 THEN o.customer_email END) AS returning_customers
FROM orders o
WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed' AND o.customer_email IS NOT NULL;

-- name: GetCustomerSummaryTop :many
SELECT
  o.customer_email,
  MAX(o.customer_name) AS name,
  COUNT(o.order_id) AS orders,
  SUM(o.amount) AS total_spent,
  CASE WHEN sc.shop_customer_id IS NOT NULL THEN true ELSE false END AS is_registered
FROM orders o
LEFT JOIN shop_customers sc ON o.shop_id = sc.shop_id AND o.customer_email = sc.email
WHERE o.shop_id = $1 AND o.created_at BETWEEN $2 AND $3 AND o.status = 'completed' AND o.customer_email IS NOT NULL
GROUP BY o.customer_email, sc.shop_customer_id
ORDER BY total_spent DESC
LIMIT 5;

-- name: GetLowStockProducts :many
SELECT
  pv.product_variation_id,
  p.title AS product_name,
  pv.sku,
  pv.description,
  pv.available_quantity AS stock
FROM product_variations pv
JOIN products p ON pv.product_id = p.product_id
WHERE pv.shop_id = $1 AND pv.available_quantity <= $2
ORDER BY pv.available_quantity ASC; 