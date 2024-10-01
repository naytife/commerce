-- Modify "product_variations" table
ALTER TABLE "product_variations" ADD CONSTRAINT "product_variations_slug_shop_id_key" UNIQUE ("slug", "shop_id");
