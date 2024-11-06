-- Modify "shops" table
ALTER TABLE "shops" ADD COLUMN "whatsapp_phone_number" character varying(16) NULL, ADD COLUMN "whatsapp_url" text NULL, ADD COLUMN "facebook_url" text NULL, ADD COLUMN "instagram_url" text NULL;
-- Drop "facebooks" table
DROP TABLE "facebooks";
-- Drop "whatsapps" table
DROP TABLE "whatsapps";
