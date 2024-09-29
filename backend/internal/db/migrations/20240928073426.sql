-- Modify "facebooks" table
ALTER TABLE "facebooks" ADD CONSTRAINT "facebooks_shop_id_key" UNIQUE ("shop_id");
-- Modify "shop_images" table
ALTER TABLE "shop_images" ADD CONSTRAINT "shop_images_shop_id_key" UNIQUE ("shop_id");
-- Modify "whatsapps" table
ALTER TABLE "whatsapps" ADD CONSTRAINT "whatsapps_shop_id_key" UNIQUE ("shop_id");
