-- Rename a column from "provider" to "auth_provider"
ALTER TABLE "users" RENAME COLUMN "provider" TO "auth_provider";
-- Rename a column from "provider_id" to "auth_provider_id"
ALTER TABLE "users" RENAME COLUMN "provider_id" TO "auth_provider_id";
-- Create "shop_customers" table
CREATE TABLE "shop_customers" ("shop_customer_id" uuid NOT NULL DEFAULT gen_random_uuid(), "shop_id" bigint NOT NULL, "email" character varying(255) NOT NULL, "name" character varying(255) NULL, "locale" character varying(255) NULL, "profile_picture" text NULL, "verified_email" boolean NULL DEFAULT false, "auth_provider" character varying(255) NULL, "auth_provider_id" character varying(255) NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "last_login" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("shop_customer_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Rename a column from "user_id" to "shop_customer_id"
ALTER TABLE "orders" RENAME COLUMN "user_id" TO "shop_customer_id";
-- Modify "orders" table
ALTER TABLE "orders" DROP CONSTRAINT "fk_user", ADD CONSTRAINT "fk_shop_customer" FOREIGN KEY ("shop_customer_id") REFERENCES "shop_customers" ("shop_customer_id") ON UPDATE NO ACTION ON DELETE CASCADE;
