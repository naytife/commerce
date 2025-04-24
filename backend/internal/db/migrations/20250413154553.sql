-- Modify "attribute_options" table
ALTER TABLE "attribute_options" DROP CONSTRAINT "attribute_options_value_attribute_id_key", ADD CONSTRAINT "attribute_options_value_attribute_id_shop_id_key" UNIQUE ("value", "attribute_id", "shop_id");
-- Modify "attributes" table
ALTER TABLE "attributes" DROP CONSTRAINT "attributes_title_product_type_id_key", ADD CONSTRAINT "attributes_title_product_type_id_shop_id_key" UNIQUE ("title", "product_type_id", "shop_id");
