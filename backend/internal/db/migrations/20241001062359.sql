-- Create "product_variations" table
CREATE TABLE "product_variations" ("product_variation_id" bigserial NOT NULL, "slug" character varying(50) NOT NULL, "description" character varying(255) NOT NULL, "price" numeric(10,2) NOT NULL, "available_quantity" bigint NOT NULL, "attributes" jsonb NOT NULL DEFAULT '{}', "seo_description" text NULL, "seo_keywords" text[] NULL, "seo_title" character varying(255) NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "product_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_variation_id"), CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "products" ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE);

-- SET RLS for product_variations
ALTER TABLE product_variations ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_variations
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);