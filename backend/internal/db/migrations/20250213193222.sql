-- Modify "product_variations" table
ALTER TABLE "product_variations" DROP COLUMN "status";
-- Modify "products" table
ALTER TABLE "products" ADD COLUMN "status" "product_status" NOT NULL DEFAULT 'DRAFT';
