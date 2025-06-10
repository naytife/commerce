-- Create enum type "payment_method_type"
CREATE TYPE "payment_method_type" AS ENUM ('flutterwave', 'paystack', 'paypal', 'stripe');
-- Create enum type "order_status_type"
CREATE TYPE "order_status_type" AS ENUM ('pending', 'processing', 'completed', 'cancelled', 'refunded');
-- Create enum type "payment_status_type"
CREATE TYPE "payment_status_type" AS ENUM ('pending', 'paid', 'failed', 'refunded', 'partial_refund');
-- Create enum type "shipping_status_type"
CREATE TYPE "shipping_status_type" AS ENUM ('pending', 'shipped', 'delivered', 'cancelled', 'returned');
-- Modify "order_items" table
ALTER TABLE "order_items" ADD COLUMN "shop_id" bigint NOT NULL, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "orders" table
ALTER TABLE "orders" ALTER COLUMN "status" TYPE "order_status_type" USING "status"::"order_status_type", ALTER COLUMN "status" SET DEFAULT 'pending', DROP COLUMN "total_price", ALTER COLUMN "user_id" DROP NOT NULL, ADD COLUMN "amount" numeric(10,2) NOT NULL, ADD COLUMN "discount" numeric(10,2) NOT NULL DEFAULT 0, ADD COLUMN "shipping_cost" numeric(10,2) NOT NULL DEFAULT 0, ADD COLUMN "tax" numeric(10,2) NOT NULL DEFAULT 0, ADD COLUMN "shipping_address" text NOT NULL, ADD COLUMN "payment_method" "payment_method_type" NOT NULL, ADD COLUMN "payment_status" "payment_status_type" NOT NULL DEFAULT 'pending', ADD COLUMN "shipping_method" character varying(10) NOT NULL, ADD COLUMN "shipping_status" "shipping_status_type" NOT NULL DEFAULT 'pending', ADD COLUMN "transaction_id" text NULL, ADD COLUMN "username" character varying(50) NOT NULL, ADD COLUMN "shop_id" bigint NOT NULL, ADD CONSTRAINT "fk_shop" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create "shop_payment_methods" table
CREATE TABLE "shop_payment_methods" ("payment_method_id" bigserial NOT NULL, "shop_id" bigint NOT NULL, "method_type" "payment_method_type" NOT NULL, "is_enabled" boolean NOT NULL DEFAULT false, "attributes" jsonb NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("payment_method_id"), CONSTRAINT "shop_payment_methods_shop_id_method_type_key" UNIQUE ("shop_id", "method_type"), CONSTRAINT "shop_payment_methods_shop_id_fkey" FOREIGN KEY ("shop_id") REFERENCES "shops" ("shop_id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Drop "shopping_cart_items" table
DROP TABLE "shopping_cart_items";
-- Drop "shopping_cart" table
DROP TABLE "shopping_cart";
