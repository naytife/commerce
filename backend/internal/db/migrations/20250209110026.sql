-- Create enum type "product_status"
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_status') THEN
        CREATE TYPE product_status AS ENUM ('DRAFT', 'PUBLISHED', 'ARCHIVED');
    END IF;
END$$;
-- Create enum type "attribute_data_type"
CREATE TYPE "attribute_data_type" AS ENUM ('Text', 'Number', 'Date', 'Option');
-- Create enum type "attribute_unit"
CREATE TYPE "attribute_unit" AS ENUM ('KG', 'GB', 'INCH');
-- Create enum type "attribute_applies_to"
CREATE TYPE "attribute_applies_to" AS ENUM ('Product', 'ProductVariation');
-- Modify "attributes" table
ALTER TABLE "attributes" DROP COLUMN "input_type", ALTER COLUMN "unit" TYPE "attribute_unit" USING CASE WHEN "unit" IN ('KG', 'GB', 'INCH') THEN "unit"::"attribute_unit" ELSE NULL END, ADD COLUMN "data_type" "attribute_data_type" NOT NULL DEFAULT 'Text', ADD COLUMN "applies_to" "attribute_applies_to" NOT NULL DEFAULT 'Product';
-- Create "attribute_options" table
CREATE TABLE "attribute_options" ("attribute_option_id" bigserial NOT NULL, "value" character varying(50) NOT NULL, "shop_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, PRIMARY KEY ("attribute_option_id"), CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "product_attribute_values" table
CREATE TABLE "product_attribute_values" ("product_attribute_value_id" bigserial NOT NULL, "value" character varying(50) NULL, "attribute_option_id" bigint NULL, "product_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_attribute_value_id"), CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES "attribute_options" ("attribute_option_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Modify "product_variations" table
ALTER TABLE "product_variations" ALTER COLUMN "status" TYPE "product_status" USING "status"::"product_status";
-- Create "product_variation_attribute_values" table
CREATE TABLE "product_variation_attribute_values" ("product_variation_attribute_value_id" bigserial NOT NULL, "value" character varying(50) NULL, "attribute_option_id" bigint NULL, "product_variation_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_variation_attribute_value_id"), CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES "attributes" ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES "attribute_options" ("attribute_option_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_product_variation" FOREIGN KEY ("product_variation_id") REFERENCES "product_variations" ("product_variation_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Drop "attribute_values" table
DROP TABLE "attribute_values";
