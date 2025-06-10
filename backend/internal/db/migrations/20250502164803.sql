-- Modify "shop_customers" table
ALTER TABLE "shop_customers" ADD COLUMN "sub" character varying(255) NULL, ADD CONSTRAINT "shop_customers_email_shop_id_key" UNIQUE ("email", "shop_id"), ADD CONSTRAINT "shop_customers_sub_key" UNIQUE ("sub"), ADD CONSTRAINT "shop_customers_sub_shop_id_key" UNIQUE ("sub", "shop_id");
