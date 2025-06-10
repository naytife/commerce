CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),     
    sub VARCHAR(255) UNIQUE,      
    email VARCHAR(255) UNIQUE,          
    auth_provider VARCHAR(255),        
    auth_provider_id VARCHAR(255),        
    name VARCHAR(255),   
    locale VARCHAR(255),                 
    profile_picture TEXT,
    verified_email BOOLEAN DEFAULT FALSE,             
    created_at TIMESTAMP DEFAULT NOW(),   
    last_login TIMESTAMP DEFAULT NOW()                 
);


CREATE TABLE shops (
    shop_id BIGSERIAL PRIMARY KEY,
    owner_id UUID NOT NULL,
    title VARCHAR(50) NOT NULL,
    subdomain VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(50) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    status VARCHAR(10) NOT NULL,
    about TEXT,
    address TEXT,
    phone_number VARCHAR(16),
    whatsapp_phone_number VARCHAR(16),
    whatsapp_link TEXT,
    facebook_link TEXT,
    instagram_link TEXT,
    seo_description TEXT,
    seo_keywords TEXT[],
    seo_title VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(user_id) 
);

CREATE TABLE shop_customers(
    shop_customer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sub VARCHAR(255) UNIQUE,
    shop_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,   
    name VARCHAR(255),
    locale VARCHAR(255),
    profile_picture TEXT,
    verified_email BOOLEAN DEFAULT FALSE, 
    auth_provider VARCHAR(255),
    auth_provider_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    last_login TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (email, shop_id),
    UNIQUE (sub, shop_id),
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE    
);

CREATE TABLE shop_images (
    shop_image_id BIGSERIAL PRIMARY KEY,
    favicon_url TEXT,
    logo_url TEXT,
    logo_url_dark TEXT,
    banner_url TEXT,
    banner_url_dark TEXT,
    cover_image_url TEXT,
    cover_image_url_dark TEXT,
    shop_id BIGINT NOT NULL UNIQUE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE categories (
    category_id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    parent_id BIGINT,
    shop_id BIGINT NOT NULL,
    UNIQUE (title, shop_id),
    UNIQUE (slug, shop_id),
    CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES categories(category_id),
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_lower_category_title_shop ON categories (LOWER(title), shop_id);

CREATE TABLE product_types (
    product_type_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    shippable BOOLEAN NOT NULL DEFAULT TRUE,
    digital BOOLEAN NOT NULL DEFAULT FALSE,
    sku_substring VARCHAR(10),
    shop_id BIGINT NOT NULL,
    UNIQUE (title, shop_id),
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TYPE product_status AS ENUM('DRAFT','PUBLISHED', 'ARCHIVED');
CREATE TABLE products(
    product_id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    status product_status NOT NULL DEFAULT 'DRAFT'::product_status,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_type_id BIGINT NOT NULL,
    category_id BIGINT,
    shop_id BIGINT NOT NULL,
    UNIQUE (slug, shop_id),
    CONSTRAINT fk_product_type FOREIGN KEY (product_type_id) REFERENCES product_types(product_type_id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_lower_product_title_shop ON products (LOWER(title), shop_id);

CREATE TABLE product_images(
    product_image_id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    alt VARCHAR(255) NOT NULL,
    product_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE product_variations(
    product_variation_id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    available_quantity BIGINT NOT NULL,
    seo_description TEXT,
    seo_keywords TEXT[],
    seo_title VARCHAR(255),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (sku, shop_id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE 
);

CREATE TYPE attribute_data_type AS ENUM('Text', 'Number', 'Date', 'Option', 'Color');
CREATE TYPE attribute_unit AS ENUM('KG', 'GB', 'INCH');
CREATE TYPE attribute_applies_to AS ENUM('Product', 'ProductVariation');
CREATE TABLE attributes(
    attribute_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    data_type attribute_data_type NOT NULL DEFAULT 'Text'::attribute_data_type,
    unit attribute_unit,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    applies_to attribute_applies_to NOT NULL DEFAULT 'Product'::attribute_applies_to,
    product_type_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (title, product_type_id, shop_id),
    CONSTRAINT fk_product_type FOREIGN KEY (product_type_id) REFERENCES product_types(product_type_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_lower_title_product_type ON attributes (LOWER(title), product_type_id);

CREATE TABLE attribute_options(
    attribute_option_id BIGSERIAL PRIMARY KEY,
    value VARCHAR(50) NOT NULL,
    shop_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    UNIQUE (value, attribute_id, shop_id),
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX unique_lower_value_attribute ON attribute_options (LOWER(value), attribute_id);

CREATE TABLE product_attribute_values(
    product_attribute_value_id BIGSERIAL PRIMARY KEY,
    value VARCHAR(50),
    attribute_option_id BIGINT,
    product_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (product_id, attribute_id, shop_id),
    CONSTRAINT fk_attribute_option FOREIGN KEY (attribute_option_id) REFERENCES attribute_options(attribute_option_id) ON DELETE SET NULL,
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE product_variation_attribute_values(
    product_variation_attribute_value_id BIGSERIAL PRIMARY KEY,
    value VARCHAR(50),
    attribute_option_id BIGINT,
    product_variation_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (product_variation_id, attribute_id, shop_id),
    CONSTRAINT fk_attribute_option FOREIGN KEY (attribute_option_id) REFERENCES attribute_options(attribute_option_id) ON DELETE SET NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TYPE payment_method_type AS ENUM (
  'flutterwave',
  'paystack',
  'paypal',
  'stripe'
);
CREATE TABLE shop_payment_methods(
  payment_method_id BIGSERIAL PRIMARY KEY,
  shop_id BIGINT NOT NULL REFERENCES shops(shop_id),
  method_type payment_method_type NOT NULL,
  is_enabled BOOLEAN NOT NULL DEFAULT false,
  attributes JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  UNIQUE (shop_id, method_type)
);

-- Order lifecycle states
CREATE TYPE order_status_type AS ENUM (
  'pending',    -- Order placed but not processed
  'processing', -- Payment confirmed, preparing for shipping
  'completed',  -- Order delivered and closed
  'cancelled',  -- Order cancelled before completion
  'refunded'    -- Full refund issued
);

-- Payment outcomes
CREATE TYPE payment_status_type AS ENUM (
  'pending',       -- Awaiting payment
  'paid',          -- Full payment received
  'failed',        -- Payment attempt failed
  'refunded',      -- Full refund processed
  'partial_refund' -- Partial refund issued
);

-- Shipping progress states
CREATE TYPE shipping_status_type AS ENUM (
  'pending',    -- Awaiting fulfillment
  'shipped',    -- Dispatched to carrier
  'delivered',  -- Received by customer
  'cancelled',  -- Shipping cancelled
  'returned'    -- Returned to seller
);

CREATE TABLE orders (
  order_id BIGSERIAL PRIMARY KEY,
  status order_status_type NOT NULL DEFAULT 'pending',
  amount DECIMAL(10, 2) NOT NULL,
  discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
  shipping_cost DECIMAL(10, 2) NOT NULL DEFAULT 0,
  tax DECIMAL(10, 2) NOT NULL DEFAULT 0,
  shipping_address TEXT NOT NULL,
  payment_method payment_method_type NOT NULL, -- Reuse existing ENUM
  payment_status payment_status_type NOT NULL DEFAULT 'pending',
  shipping_method VARCHAR(10) NOT NULL,
  shipping_status shipping_status_type NOT NULL DEFAULT 'pending',
  transaction_id TEXT,
  username VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  shop_customer_id UUID,
  shop_id BIGINT NOT NULL,
  customer_name VARCHAR(100) NOT NULL,
  customer_email VARCHAR(100),
  customer_phone VARCHAR(50),
  CONSTRAINT fk_shop_customer FOREIGN KEY (shop_customer_id) REFERENCES shop_customers(shop_customer_id) ON DELETE CASCADE,
  CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE order_items(
    order_item_id BIGSERIAL PRIMARY KEY,
    quantity BIGINT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_variation_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

-- Stock movements table for inventory tracking
CREATE TABLE stock_movements (
    movement_id BIGSERIAL PRIMARY KEY,
    product_variation_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    movement_type VARCHAR(50) NOT NULL, -- 'SALE', 'RESTOCK', 'ADJUSTMENT', 'RETURN'
    quantity_change INT NOT NULL, -- positive for increase, negative for decrease
    quantity_before INT NOT NULL,
    quantity_after INT NOT NULL,
    reference_id BIGINT, -- order_id for sales, adjustment_id for manual adjustments
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

-- SET RLS for categories
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON categories
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for products
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON products
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for shop_images
ALTER TABLE shop_images ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON shop_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_images
ALTER TABLE product_images ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_images
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_variations
ALTER TABLE product_variations ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_variations
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for attributes
ALTER TABLE attributes ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON attributes
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for attribute_options
ALTER TABLE attribute_options ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON attribute_options
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_attribute_values
ALTER TABLE product_attribute_values ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_attribute_values
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for product_variation_attribute_values
ALTER TABLE product_variation_attribute_values ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON product_variation_attribute_values
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for orders
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON orders
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for order_items
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON order_items
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);

-- SET RLS for stock_movements
ALTER TABLE stock_movements ENABLE ROW LEVEL SECURITY;

CREATE POLICY shop_policy ON stock_movements
FOR ALL
USING (shop_id = current_setting('commerce.current_shop_id')::int)
WITH CHECK (shop_id = current_setting('commerce.current_shop_id')::int);
