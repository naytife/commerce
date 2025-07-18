CREATE SCHEMA IF NOT EXISTS naytife_schema;
-- Create enum type "product_status"
CREATE TYPE naytife_schema.product_status AS ENUM ('DRAFT', 'PUBLISHED', 'ARCHIVED');
-- Create enum type "attribute_data_type"
CREATE TYPE naytife_schema.attribute_data_type AS ENUM ('Text', 'Number', 'Date', 'Option', 'Color');
-- Create enum type "attribute_unit"
CREATE TYPE naytife_schema.attribute_unit AS ENUM ('KG', 'GB', 'INCH');
-- Create enum type "attribute_applies_to"
CREATE TYPE naytife_schema.attribute_applies_to AS ENUM ('Product', 'ProductVariation');
-- Create enum type "payment_method_type"
CREATE TYPE naytife_schema.payment_method_type AS ENUM ('flutterwave', 'paystack', 'paypal', 'stripe', 'pay_on_delivery');
-- Create enum type "order_status_type"
CREATE TYPE naytife_schema.order_status_type AS ENUM ('pending', 'processing', 'completed', 'cancelled', 'refunded');
-- Create enum type "payment_status_type"
CREATE TYPE naytife_schema.payment_status_type AS ENUM ('pending', 'paid', 'failed', 'refunded', 'partial_refund');
-- Create enum type "shipping_status_type"
CREATE TYPE naytife_schema.shipping_status_type AS ENUM ('pending', 'shipped', 'delivered', 'cancelled', 'returned');
-- Create "attribute_options" table
CREATE TABLE naytife_schema.attribute_options ("attribute_option_id" bigserial NOT NULL, "value" character varying(50) NOT NULL, "shop_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, PRIMARY KEY ("attribute_option_id"), CONSTRAINT "attribute_options_value_attribute_id_shop_id_key" UNIQUE ("value", "attribute_id", "shop_id"));
-- Create index "unique_lower_value_attribute" to table: "attribute_options"
CREATE UNIQUE INDEX unique_lower_value_attribute ON naytife_schema.attribute_options ((lower((value)::text)), "attribute_id");
-- Create "attributes" table
CREATE TABLE naytife_schema.attributes ("attribute_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "data_type" naytife_schema.attribute_data_type NOT NULL DEFAULT 'Text', "unit" naytife_schema.attribute_unit NULL, "required" boolean NOT NULL DEFAULT false, "applies_to" naytife_schema.attribute_applies_to NOT NULL DEFAULT 'Product', "product_type_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("attribute_id"), CONSTRAINT "attributes_title_product_type_id_shop_id_key" UNIQUE ("title", "product_type_id", "shop_id"));
-- Create index "unique_lower_title_product_type" to table: "attributes"
CREATE UNIQUE INDEX unique_lower_title_product_type ON naytife_schema.attributes ((lower((title)::text)), "product_type_id");
-- Create "categories" table
CREATE TABLE naytife_schema.categories ("category_id" bigserial NOT NULL, "slug" character varying(50) NOT NULL, "title" character varying(50) NOT NULL, "description" character varying(255) NULL, "parent_id" bigint NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("category_id"), CONSTRAINT "categories_slug_shop_id_key" UNIQUE ("slug", "shop_id"), CONSTRAINT "categories_title_shop_id_key" UNIQUE ("title", "shop_id"), CONSTRAINT "fk_parent" FOREIGN KEY ("parent_id") REFERENCES naytife_schema.categories ("category_id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "unique_lower_category_title_shop" to table: "categories"
CREATE UNIQUE INDEX unique_lower_category_title_shop ON naytife_schema.categories ((lower((title)::text)), "shop_id");
-- Create "order_items" table
CREATE TABLE naytife_schema.order_items ("order_item_id" bigserial NOT NULL, "quantity" bigint NOT NULL, "price" numeric(10,2) NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "product_variation_id" bigint NOT NULL, "order_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("order_item_id"));
-- Create "orders" table
CREATE TABLE naytife_schema.orders ("order_id" bigserial NOT NULL, "status" naytife_schema.order_status_type NOT NULL DEFAULT 'pending', "amount" numeric(10,2) NOT NULL, "discount" numeric(10,2) NOT NULL DEFAULT 0, "shipping_cost" numeric(10,2) NOT NULL DEFAULT 0, "tax" numeric(10,2) NOT NULL DEFAULT 0, "shipping_address" text NOT NULL, "payment_method" naytife_schema.payment_method_type NOT NULL, "payment_status" naytife_schema.payment_status_type NOT NULL DEFAULT 'pending', "shipping_method" character varying(10) NOT NULL, "shipping_status" naytife_schema.shipping_status_type NOT NULL DEFAULT 'pending', "transaction_id" text NULL, "username" character varying(50) NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "shop_customer_id" uuid NULL, "shop_id" bigint NOT NULL, "customer_name" character varying(100) NOT NULL, "customer_email" character varying(100) NULL, "customer_phone" character varying(50) NULL, PRIMARY KEY ("order_id"));
-- Create "product_attribute_values" table
CREATE TABLE naytife_schema.product_attribute_values ("product_attribute_value_id" bigserial NOT NULL, "value" character varying(50) NULL, "attribute_option_id" bigint NULL, "product_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_attribute_value_id"), CONSTRAINT "product_attribute_values_product_id_attribute_id_shop_id_key" UNIQUE ("product_id", "attribute_id", "shop_id"));
-- Create "product_images" table
CREATE TABLE naytife_schema.product_images ("product_image_id" bigserial NOT NULL, "url" text NOT NULL, "alt" character varying(255) NOT NULL, "product_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_image_id"));
-- Create "product_types" table
CREATE TABLE naytife_schema.product_types ("product_type_id" bigserial NOT NULL, "title" character varying(50) NOT NULL, "shippable" boolean NOT NULL DEFAULT true, "digital" boolean NOT NULL DEFAULT false, "sku_substring" character varying(10) NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_type_id"), CONSTRAINT "product_types_title_shop_id_key" UNIQUE ("title", "shop_id"));
-- Create "product_variation_attribute_values" table
CREATE TABLE naytife_schema.product_variation_attribute_values ("product_variation_attribute_value_id" bigserial NOT NULL, "value" character varying(50) NULL, "attribute_option_id" bigint NULL, "product_variation_id" bigint NOT NULL, "attribute_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_variation_attribute_value_id"), CONSTRAINT "product_variation_attribute_v_product_variation_id_attribut_key" UNIQUE ("product_variation_id", "attribute_id", "shop_id"));
-- Create "product_variations" table
CREATE TABLE naytife_schema.product_variations ("product_variation_id" bigserial NOT NULL, "sku" character varying(50) NOT NULL, "description" character varying(255) NOT NULL, "price" numeric(10,2) NOT NULL, "available_quantity" bigint NOT NULL, "seo_description" text NULL, "seo_keywords" text[] NULL, "seo_title" character varying(255) NULL, "is_default" boolean NOT NULL DEFAULT false, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "product_id" bigint NOT NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_variation_id"), CONSTRAINT "product_variations_sku_shop_id_key" UNIQUE ("sku", "shop_id"));
-- Create "products" table
CREATE TABLE naytife_schema.products ("product_id" bigserial NOT NULL, "slug" character varying(50) NOT NULL, "title" character varying(50) NOT NULL, "description" character varying(255) NOT NULL, "status" naytife_schema.product_status NOT NULL DEFAULT 'DRAFT', "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "product_type_id" bigint NOT NULL, "category_id" bigint NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("product_id"), CONSTRAINT "products_slug_shop_id_key" UNIQUE ("slug", "shop_id"));
-- Create index "unique_lower_product_title_shop" to table: "products"
CREATE UNIQUE INDEX unique_lower_product_title_shop ON naytife_schema.products ((lower((title)::text)), "shop_id");
-- Create "shop_customers" table
CREATE TABLE naytife_schema.shop_customers ("shop_customer_id" uuid NOT NULL DEFAULT gen_random_uuid(), "sub" character varying(255) NULL, "shop_id" bigint NOT NULL, "email" character varying(255) NOT NULL, "name" character varying(255) NULL, "locale" character varying(255) NULL, "profile_picture" text NULL, "verified_email" boolean NULL DEFAULT false, "auth_provider" character varying(255) NULL, "auth_provider_id" character varying(255) NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "last_login" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("shop_customer_id"), CONSTRAINT "shop_customers_email_shop_id_key" UNIQUE ("email", "shop_id"), CONSTRAINT "shop_customers_sub_key" UNIQUE ("sub"), CONSTRAINT "shop_customers_sub_shop_id_key" UNIQUE ("sub", "shop_id"));
-- Create "shop_data_updates" table
CREATE TABLE naytife_schema.shop_data_updates ("update_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "data_type" character varying(50) NOT NULL, "status" character varying(50) NOT NULL DEFAULT 'updating', "changes_summary" jsonb NULL, "started_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed_at" timestamptz NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("update_id"));
-- Create index "idx_shop_data_updates_shop_id" to table: "shop_data_updates"
CREATE INDEX idx_shop_data_updates_shop_id ON naytife_schema.shop_data_updates ("shop_id");
-- Create index "idx_shop_data_updates_status" to table: "shop_data_updates"
CREATE INDEX idx_shop_data_updates_status ON naytife_schema.shop_data_updates ("status");
-- Create "shop_deployment_urls" table
CREATE TABLE naytife_schema.shop_deployment_urls ("url_id" bigserial NOT NULL, "deployment_id" bigint NOT NULL, "url_type" character varying(50) NOT NULL, "url" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("url_id"));
-- Create "shop_deployments" table
CREATE TABLE naytife_schema.shop_deployments ("deployment_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "template_name" character varying(100) NOT NULL, "template_version" character varying(50) NOT NULL DEFAULT 'latest', "status" character varying(50) NOT NULL DEFAULT 'deploying', "deployment_type" character varying(50) NOT NULL DEFAULT 'full', "message" text NULL, "started_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "completed_at" timestamptz NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("deployment_id"));
-- Create index "idx_shop_deployments_shop_id" to table: "shop_deployments"
CREATE INDEX idx_shop_deployments_shop_id ON naytife_schema.shop_deployments ("shop_id");
-- Create index "idx_shop_deployments_status" to table: "shop_deployments"
CREATE INDEX idx_shop_deployments_status ON naytife_schema.shop_deployments ("status");
-- Create index "idx_shop_deployments_template" to table: "shop_deployments"
CREATE INDEX idx_shop_deployments_template ON naytife_schema.shop_deployments ("template_name");
-- Create "shop_images" table
CREATE TABLE naytife_schema.shop_images ("shop_image_id" bigserial NOT NULL, "favicon_url" text NULL, "logo_url" text NULL, "logo_url_dark" text NULL, "banner_url" text NULL, "banner_url_dark" text NULL, "cover_image_url" text NULL, "cover_image_url_dark" text NULL, "shop_id" bigint NOT NULL, PRIMARY KEY ("shop_image_id"), CONSTRAINT "shop_images_shop_id_key" UNIQUE ("shop_id"));
-- Create "shop_payment_methods" table
CREATE TABLE naytife_schema.shop_payment_methods ("payment_method_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "method_type" naytife_schema.payment_method_type NOT NULL, "is_enabled" boolean NOT NULL DEFAULT false, "attributes" jsonb NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("payment_method_id"), CONSTRAINT "shop_payment_methods_shop_id_method_type_key" UNIQUE ("shop_id", "method_type"));
-- Create "shops" table
CREATE TABLE naytife_schema.shops ("shop_id" bigserial NOT NULL, "owner_id" uuid NOT NULL, "title" character varying(50) NOT NULL, "subdomain" character varying(50) NOT NULL, "email" character varying(50) NOT NULL, "currency_code" character varying(3) NOT NULL, "status" character varying(10) NOT NULL, "about" text NULL, "address" text NULL, "phone_number" character varying(16) NULL, "whatsapp_phone_number" character varying(16) NULL, "whatsapp_link" text NULL, "facebook_link" text NULL, "instagram_link" text NULL, "seo_description" text NULL, "seo_keywords" text[] NULL, "seo_title" character varying(255) NULL, "current_template" character varying(100) NULL, "last_deployment_id" bigint NULL, "last_data_update_at" timestamptz NULL, "updated_at" timestamptz NOT NULL DEFAULT now(), "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("shop_id"), CONSTRAINT "shops_subdomain_key" UNIQUE ("subdomain"));
-- Create "stock_movements" table
CREATE TABLE naytife_schema.stock_movements ("movement_id" bigserial NOT NULL, "product_variation_id" bigint NOT NULL, "shop_id" bigint NOT NULL, "movement_type" character varying(50) NOT NULL, "quantity_change" integer NOT NULL, "quantity_before" integer NOT NULL, "quantity_after" integer NOT NULL, "reference_id" bigint NULL, "notes" text NULL, "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("movement_id"));
-- Create "users" table
CREATE TABLE naytife_schema.users ("user_id" uuid NOT NULL DEFAULT gen_random_uuid(), "sub" character varying(255) NULL, "email" character varying(255) NULL, "auth_provider" character varying(255) NULL, "auth_provider_id" character varying(255) NULL, "name" character varying(255) NULL, "locale" character varying(255) NULL, "profile_picture" text NULL, "verified_email" boolean NULL DEFAULT false, "created_at" timestamp NULL DEFAULT now(), "last_login" timestamp NULL DEFAULT now(), PRIMARY KEY ("user_id"), CONSTRAINT "users_email_key" UNIQUE ("email"), CONSTRAINT "users_sub_key" UNIQUE ("sub"));
-- Modify "attribute_options" table
ALTER TABLE naytife_schema.attribute_options ADD CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES naytife_schema.attributes ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "attributes" table
ALTER TABLE naytife_schema.attributes ADD CONSTRAINT "fk_product_type" FOREIGN KEY ("product_type_id") REFERENCES naytife_schema.product_types ("product_type_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "categories" table
ALTER TABLE naytife_schema.categories ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "order_items" table
ALTER TABLE naytife_schema.order_items ADD CONSTRAINT "fk_order" FOREIGN KEY ("order_id") REFERENCES naytife_schema.orders ("order_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_product_variation" FOREIGN KEY ("product_variation_id") REFERENCES naytife_schema.product_variations ("product_variation_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "orders" table
ALTER TABLE naytife_schema.orders ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop_customer" FOREIGN KEY ("shop_customer_id") REFERENCES naytife_schema.shop_customers ("shop_customer_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "product_attribute_values" table
ALTER TABLE naytife_schema.product_attribute_values ADD CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES naytife_schema.attributes ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES naytife_schema.attribute_options ("attribute_option_id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES naytife_schema.products ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "product_images" table
ALTER TABLE naytife_schema.product_images ADD CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES naytife_schema.products ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "product_types" table
ALTER TABLE naytife_schema.product_types ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "product_variation_attribute_values" table
ALTER TABLE naytife_schema.product_variation_attribute_values ADD CONSTRAINT "fk_attribute" FOREIGN KEY ("attribute_id") REFERENCES naytife_schema.attributes ("attribute_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_attribute_option" FOREIGN KEY ("attribute_option_id") REFERENCES naytife_schema.attribute_options ("attribute_option_id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_product_variation" FOREIGN KEY ("product_variation_id") REFERENCES naytife_schema.product_variations ("product_variation_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "product_variations" table
ALTER TABLE naytife_schema.product_variations ADD CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES naytife_schema.products ("product_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "products" table
ALTER TABLE naytife_schema.products ADD CONSTRAINT "fk_category" FOREIGN KEY ("category_id") REFERENCES naytife_schema.categories ("category_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_product_type" FOREIGN KEY ("product_type_id") REFERENCES naytife_schema.product_types ("product_type_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_customers" table
ALTER TABLE naytife_schema.shop_customers ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_data_updates" table
ALTER TABLE naytife_schema.shop_data_updates ADD CONSTRAINT "shop_data_updates_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_deployment_urls" table
ALTER TABLE naytife_schema.shop_deployment_urls ADD CONSTRAINT "shop_deployment_urls_deployment_id_fkey" FOREIGN KEY ("deployment_id") REFERENCES naytife_schema.shop_deployments ("deployment_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_deployments" table
ALTER TABLE naytife_schema.shop_deployments ADD CONSTRAINT "shop_deployments_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_images" table
ALTER TABLE naytife_schema.shop_images ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "shop_payment_methods" table
ALTER TABLE naytife_schema.shop_payment_methods ADD CONSTRAINT "shop_payment_methods_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "shops" table
ALTER TABLE naytife_schema.shops ADD CONSTRAINT "fk_last_deployment" FOREIGN KEY ("last_deployment_id") REFERENCES naytife_schema.shop_deployments ("deployment_id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_owner" FOREIGN KEY ("owner_id") REFERENCES naytife_schema.users ("user_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "stock_movements" table
ALTER TABLE naytife_schema.stock_movements ADD CONSTRAINT "fk_product_variation" FOREIGN KEY ("product_variation_id") REFERENCES naytife_schema.product_variations ("product_variation_id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES naytife_schema.shops ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;

-- Enable RLS and create policies for multi-tenant security

-- SET RLS for categories
ALTER TABLE naytife_schema.categories ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.categories
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for products
ALTER TABLE naytife_schema.products ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.products
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shop_images
ALTER TABLE naytife_schema.shop_images ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_images
ALTER TABLE naytife_schema.product_images ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.product_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_variations
ALTER TABLE naytife_schema.product_variations ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.product_variations
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for attributes
ALTER TABLE naytife_schema.attributes ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.attributes
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for attribute_options
ALTER TABLE naytife_schema.attribute_options ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.attribute_options
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_attribute_values
ALTER TABLE naytife_schema.product_attribute_values ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.product_attribute_values
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_variation_attribute_values
ALTER TABLE naytife_schema.product_variation_attribute_values ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.product_variation_attribute_values
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for orders
ALTER TABLE naytife_schema.orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.orders
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for order_items
ALTER TABLE naytife_schema.order_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.order_items
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for stock_movements
ALTER TABLE naytife_schema.stock_movements ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.stock_movements
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shop_deployments
ALTER TABLE naytife_schema.shop_deployments ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_deployments
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shop_deployment_urls
ALTER TABLE naytife_schema.shop_deployment_urls ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_deployment_urls
FOR ALL
USING (deployment_id IN (SELECT deployment_id FROM naytife_schema.shop_deployments WHERE shop_id = current_setting('commerce.current_shop_id')::int));

-- SET RLS for shop_data_updates
ALTER TABLE naytife_schema.shop_data_updates ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_data_updates
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shops
ALTER TABLE naytife_schema.shops ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shops
FOR ALL
USING (owner_id = current_setting('commerce.current_owner_id')::uuid)
WITH CHECK (owner_id = current_setting('commerce.current_owner_id')::uuid);

-- SET RLS for shop_customers
ALTER TABLE naytife_schema.shop_customers ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_customers
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_types
ALTER TABLE naytife_schema.product_types ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.product_types
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shop_payment_methods
ALTER TABLE naytife_schema.shop_payment_methods ENABLE ROW LEVEL SECURITY;
CREATE POLICY shop_policy ON naytife_schema.shop_payment_methods
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);
