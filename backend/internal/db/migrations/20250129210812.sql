-- Modify "product_variations" table
ALTER TABLE "product_variations" ALTER COLUMN "sku" DROP DEFAULT;
-- Modify "products" table
ALTER TABLE "products" ALTER COLUMN "product_type_id" SET NOT NULL;
