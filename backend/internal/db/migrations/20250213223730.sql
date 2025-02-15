-- Modify "product_attribute_values" table
ALTER TABLE "product_attribute_values" ADD CONSTRAINT "product_attribute_values_product_id_attribute_id_shop_id_key" UNIQUE ("product_id", "attribute_id", "shop_id");
-- Modify "product_variation_attribute_values" table
ALTER TABLE "product_variation_attribute_values" ADD CONSTRAINT "product_variation_attribute_v_product_variation_id_attribut_key" UNIQUE ("product_variation_id", "attribute_id", "shop_id");
