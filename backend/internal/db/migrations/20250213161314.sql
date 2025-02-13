-- Modify "attribute_options" table
ALTER TABLE "attribute_options" ADD CONSTRAINT "attribute_options_value_attribute_id_key" UNIQUE ("value", "attribute_id");
-- Modify "attributes" table
ALTER TABLE "attributes" ADD CONSTRAINT "attributes_title_product_type_id_key" UNIQUE ("title", "product_type_id");
-- Modify "product_attribute_values" table
ALTER TABLE "product_attribute_values" DROP CONSTRAINT "fk_attribute_option", ADD CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES "attribute_options" ("attribute_option_id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "product_variation_attribute_values" table
ALTER TABLE "product_variation_attribute_values" DROP CONSTRAINT "fk_attribute_option", ADD CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES "attribute_options" ("attribute_option_id") ON UPDATE NO ACTION ON DELETE SET NULL;
