-- Modify "categories" table
ALTER TABLE "categories" DROP COLUMN "allowed_attributes", ADD COLUMN "category_attributes" jsonb NOT NULL DEFAULT '{}';
