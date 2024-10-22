-- Modify "users" table
ALTER TABLE "users" DROP CONSTRAINT "users_auth0_sub_key", ADD COLUMN "provider_id" character varying(255) NULL, ADD COLUMN "locale" character varying(255) NULL, ADD CONSTRAINT "users_provider_id_key" UNIQUE ("provider_id");
-- Rename a column from "auth0_sub" to "provider"
ALTER TABLE "users" RENAME COLUMN "auth0_sub" TO "provider";
-- Rename a column from "profile_picture_url" to "profile_picture"
ALTER TABLE "users" RENAME COLUMN "profile_picture_url" TO "profile_picture";
