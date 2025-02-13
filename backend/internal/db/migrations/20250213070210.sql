-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "last_login" SET DEFAULT now();
-- Modify "attributes" table
ALTER TABLE "attributes" ADD COLUMN "product_type_id" bigint NOT NULL, ADD CONSTRAINT "fk_product_type" FOREIGN KEY ("product_type_id") REFERENCES "product_types" ("product_type_id") ON UPDATE NO ACTION ON DELETE CASCADE;
