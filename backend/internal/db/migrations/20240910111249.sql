-- Create "whatsapps" table
CREATE TABLE "whatsapps" ("whatsapp_id" bigserial NOT NULL, "shop_id" uuid NOT NULL, "number" character varying(20) NOT NULL, "country_code" character varying(5) NOT NULL, "url" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("whatsapp_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE NO ACTION);
