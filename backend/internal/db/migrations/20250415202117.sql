-- Modify "product_variations" table
ALTER TABLE "product_variations" DROP COLUMN "slug";
-- Modify "products" table
ALTER TABLE "products" DROP CONSTRAINT "products_title_shop_id_key", ADD COLUMN "slug" character varying(50) NOT NULL, ADD CONSTRAINT "products_slug_shop_id_key" UNIQUE ("slug", "shop_id");
