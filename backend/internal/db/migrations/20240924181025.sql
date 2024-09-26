-- SET RLS for categories
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON categories
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for products
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON products
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);