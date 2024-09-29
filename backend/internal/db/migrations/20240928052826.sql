-- Modify "whatsapps" table
ALTER TABLE "whatsapps" DROP COLUMN "country_code", DROP COLUMN "created_at";
-- Create "facebooks" table
CREATE TABLE "facebooks" ("facebook_id" bigserial NOT NULL, "handle" character varying(50) NOT NULL, "url" text NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("facebook_id"), CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);

-- SET RLS for whatsapps
ALTER TABLE whatsapps ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON whatsapps
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for facebooks
ALTER TABLE facebooks ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON facebooks
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);