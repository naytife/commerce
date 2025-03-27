-- Create index "unique_lower_value_attribute" to table: "attribute_options"
CREATE UNIQUE INDEX "unique_lower_value_attribute" ON "attribute_options" ((lower((value)::text)), "attribute_id");
-- Create index "unique_lower_title_product_type" to table: "attributes"
CREATE UNIQUE INDEX "unique_lower_title_product_type" ON "attributes" ((lower((title)::text)), "product_type_id");
-- Create index "unique_lower_category_title_shop" to table: "categories"
CREATE UNIQUE INDEX "unique_lower_category_title_shop" ON "categories" ((lower((title)::text)), "shop_id");
-- Create index "unique_lower_product_title_shop" to table: "products"
CREATE UNIQUE INDEX "unique_lower_product_title_shop" ON "products" ((lower((title)::text)), "shop_id");
-- Modify "shops" table
ALTER TABLE "shops" ADD COLUMN "subdomain" character varying(50) NULL;
