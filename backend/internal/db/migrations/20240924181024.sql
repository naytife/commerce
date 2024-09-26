-- Rename a column from "default_domain" to "domain"
ALTER TABLE "shops" RENAME COLUMN "default_domain" TO "domain";
-- Modify "shops" table
ALTER TABLE "shops" DROP CONSTRAINT "shops_default_domain_check", DROP CONSTRAINT "shops_default_domain_key", ADD CONSTRAINT "shops_domain_check" CHECK ((domain)::text ~~ '%.%'::text), ADD CONSTRAINT "shops_domain_key" UNIQUE ("domain");
-- Modify "categories" table
ALTER TABLE "categories" DROP CONSTRAINT "fk_shop", ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create "products" table
CREATE TABLE "products" ("product_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "description" character varying(255) NOT NULL, "allowed_attributes" jsonb NOT NULL DEFAULT '{}', "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "status" character varying(10) NOT NULL, "category_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_id"), CONSTRAINT "products_title_shop_id_key" UNIQUE ("title", "shop_id"), CONSTRAINT "fk_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Modify "whatsapps" table
ALTER TABLE "whatsapps" DROP CONSTRAINT "fk_shop", ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
