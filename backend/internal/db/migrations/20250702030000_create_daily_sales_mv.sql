-- Create materialized view for daily sales analytics
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_sales AS
SELECT
  shop_id,
  DATE(created_at) AS day,
  COUNT(*) AS total_orders,
  SUM(amount) AS revenue
FROM orders
WHERE status = 'completed'
GROUP BY shop_id, day;

CREATE UNIQUE INDEX IF NOT EXISTS idx_daily_sales_shop_day ON daily_sales(shop_id, day); 