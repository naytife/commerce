-- Modify "shops" table
ALTER TABLE "shops" ADD COLUMN "current_template" character varying(100) NULL, ADD COLUMN "last_deployment_id" bigint NULL, ADD COLUMN "last_data_update_at" timestamptz NULL;
-- Create "shop_data_updates" table
CREATE TABLE "shop_data_updates" ("update_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "data_type" character varying(50) NOT NULL, "status" character varying(50) NOT NULL DEFAULT 'updating', "changes_summary" jsonb NULL, "started_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed_at" timestamptz NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("update_id"));
-- Create index "idx_shop_data_updates_shop_id" to table: "shop_data_updates"
CREATE INDEX "idx_shop_data_updates_shop_id" ON "shop_data_updates" ("shop_id");
-- Create index "idx_shop_data_updates_status" to table: "shop_data_updates"
CREATE INDEX "idx_shop_data_updates_status" ON "shop_data_updates" ("status");
-- Create "shop_deployment_urls" table
CREATE TABLE "shop_deployment_urls" ("url_id" bigserial NOT NULL, "deployment_id" bigint NOT NULL, "url_type" character varying(50) NOT NULL, "url" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("url_id"));
-- Create "shop_deployments" table
CREATE TABLE "shop_deployments" ("deployment_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "template_name" character varying(100) NOT NULL, "template_version" character varying(50) NOT NULL DEFAULT 'latest', "status" character varying(50) NOT NULL DEFAULT 'deploying', "deployment_type" character varying(50) NOT NULL DEFAULT 'full', "message" text NULL, "started_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed_at" timestamptz NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("deployment_id"));
-- Create index "idx_shop_deployments_shop_id" to table: "shop_deployments"
CREATE INDEX "idx_shop_deployments_shop_id" ON "shop_deployments" ("shop_id");
-- Create index "idx_shop_deployments_status" to table: "shop_deployments"
CREATE INDEX "idx_shop_deployments_status" ON "shop_deployments" ("status");
-- Create index "idx_shop_deployments_template" to table: "shop_deployments"
CREATE INDEX "idx_shop_deployments_template" ON "shop_deployments" ("template_name");
-- Create "stock_movements" table
CREATE TABLE "stock_movements" ("movement_id" bigserial NOT NULL, "product_variation_id" bigint NOT NULL, "shop_id" bigint NOT NULL, "movement_type" character varying(50) NOT NULL, "quantity_change" integer NOT NULL, "quantity_before" integer NOT NULL, "quantity_after" integer NOT NULL, "reference_id" bigint NULL, "notes" text NULL, "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("movement_id"));
-- Modify "shops" table
ALTER TABLE "shops" ADD CONSTRAINT "fk_last_deployment" FOREIGN KEY ("last_deployment_id") REFERENCES "shop_deployments" ("deployment_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "shop_data_updates" table
ALTER TABLE "shop_data_updates" ADD CONSTRAINT "shop_data_updates_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_deployment_urls" table
ALTER TABLE "shop_deployment_urls" ADD CONSTRAINT "shop_deployment_urls_deployment_id_fkey" FOREIGN KEY ("deployment_id") REFERENCES "shop_deployments" ("deployment_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_deployments" table
ALTER TABLE "shop_deployments" ADD CONSTRAINT "shop_deployments_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "stock_movements" table
ALTER TABLE "stock_movements" ADD CONSTRAINT "fk_product_variation" FOREIGN KEY ("product_variation_id") REFERENCES "product_variations" ("product_variation_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
