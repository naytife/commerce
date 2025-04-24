-- Modify "product_variations" table
ALTER TABLE "product_variations" ADD COLUMN "is_default" boolean NOT NULL DEFAULT false;
-- Modify "shop_images" table
ALTER TABLE "shop_images" ADD COLUMN "logo_url_dark" text NULL, ADD COLUMN "banner_url_dark" text NULL, ADD COLUMN "cover_image_url_dark" text NULL;
