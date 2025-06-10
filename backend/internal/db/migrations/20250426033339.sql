-- Modify "orders" table
ALTER TABLE "orders" ADD COLUMN "customer_name" character varying(100) NOT NULL, ADD COLUMN "customer_email" character varying(100) NULL, ADD COLUMN "customer_phone" character varying(50) NULL;
