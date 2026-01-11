-- Drop the unique index on (LOWER(title), shop_id) to allow duplicate product titles within a shop
-- Slug uniqueness with product ID suffix ensures URL uniqueness instead
DROP INDEX IF EXISTS unique_lower_product_title_shop;
