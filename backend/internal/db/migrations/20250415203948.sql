-- Modify "product_variations" table
ALTER TABLE "product_variations" ADD CONSTRAINT "product_variations_sku_shop_id_key" UNIQUE ("sku", "shop_id");
