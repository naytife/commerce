-- Modify "shops" table
ALTER TABLE "shops" ADD COLUMN "owner_id" uuid NOT NULL, ADD COLUMN "status" character varying(10) NOT NULL, ADD COLUMN "seo_description" text NULL, ADD COLUMN "seo_keywords" text NULL, ADD COLUMN "seo_title" character varying(255) NULL;
