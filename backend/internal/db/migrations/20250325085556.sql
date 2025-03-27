-- Modify "shops" table
ALTER TABLE "shops" ALTER COLUMN "subdomain" SET NOT NULL, ADD CONSTRAINT "shops_subdomain_key" UNIQUE ("subdomain");
