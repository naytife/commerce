-- Modify "categories" table
ALTER TABLE "categories" DROP COLUMN "created_at", DROP COLUMN "updated_at", DROP COLUMN "category_attributes";
-- Modify "product_variations" table
ALTER TABLE "product_variations" DROP COLUMN "attributes", ADD COLUMN "sku" character varying(50) NOT NULL DEFAULT gen_random_uuid(), ADD COLUMN "status" character varying(10) NOT NULL DEFAULT 'DRAFT';
-- Modify "shopping_cart" table
ALTER TABLE "shopping_cart" DROP COLUMN "created_at", DROP COLUMN "updated_at";
-- Create "attributes" table
CREATE TABLE "attributes" ("attribute_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "input_type" character varying(50) NOT NULL, "unit" character varying(50) NULL, "required" boolean NOT NULL DEFAULT false, "shop_id" bigint NOT NULL, PRIMARY KEY ("attribute_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "attribute_values" table
CREATE TABLE "attribute_values" ("attribute_value_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "value" character varying(50) NOT NULL, "boolean_value" boolean NULL, "shop_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, PRIMARY KEY ("attribute_value_id"), CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "product_types" table
CREATE TABLE "product_types" ("product_type_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "shippable" boolean NOT NULL DEFAULT true, "digital" boolean NOT NULL DEFAULT false, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_type_id"), CONSTRAINT "product_types_title_shop_id_key" UNIQUE ("title", "shop_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Modify "products" table
ALTER TABLE "products" DROP COLUMN "allowed_attributes", DROP COLUMN "status", ADD COLUMN "product_type_id" bigint NULL, ADD CONSTRAINT "fk_product_type" FOREIGN KEY ("product_type_id") REFERENCES "product_types" ("product_type_id") ON UPDATE NO ACTION ON DELETE CASCADE;
