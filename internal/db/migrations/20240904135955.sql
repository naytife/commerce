-- Modify "shops" table
ALTER TABLE "shops" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
