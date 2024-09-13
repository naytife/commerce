-- Modify "categories" table
ALTER TABLE "categories" ADD CONSTRAINT "categories_title_shop_id_key" UNIQUE ("title", "shop_id");
