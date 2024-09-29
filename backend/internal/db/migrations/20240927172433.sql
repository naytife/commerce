-- Modify "whatsapps" table
ALTER TABLE "whatsapps" ALTER COLUMN "phone_number" TYPE character varying(16);
-- Modify "shops" table
ALTER TABLE "shops" DROP COLUMN "favicon_url", DROP COLUMN "logo_url", ALTER COLUMN "phone_number" TYPE character varying(16);
-- Create "product_images" table
CREATE TABLE "product_images" ("product_image_id" bigserial NOT NULL, "url" text NOT NULL, "alt" character varying(255) NOT NULL, "product_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_image_id"), CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "shop_images" table
CREATE TABLE "shop_images" ("shop_image_id" bigserial NOT NULL, "favicon_url" text NULL, "logo_url" text NULL, "banner_url" text NULL, "cover_image_url" text NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("shop_image_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
