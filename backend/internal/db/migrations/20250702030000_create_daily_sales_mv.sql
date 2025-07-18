CREATE SCHEMA IF NOT EXISTS naytife_schema;
-- Create materialized view for daily sales analytics
CREATE MATERIALIZED VIEW IF NOT EXISTS naytife_schema.daily_sales AS
SELECT
  shop_id,
  DATE(created_at) AS day,
  COUNT(*) AS total_orders,
  SUM(amount) AS revenue
FROM naytife_schema.orders
WHERE status = 'completed'
GROUP BY shop_id, day;

CREATE UNIQUE INDEX IF NOT EXISTS idx_daily_sales_shop_day ON naytife_schema.daily_sales(shop_id, day); 