-- Modify "users" table
ALTER TABLE "users" DROP CONSTRAINT "users_provider_id_key", ALTER COLUMN "email" DROP NOT NULL, ADD CONSTRAINT "users_email_key" UNIQUE ("email");
