-- SET RLS for shop_images
ALTER TABLE shop_images ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON shop_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_images
ALTER TABLE product_images ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);