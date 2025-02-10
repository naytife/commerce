CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),           
    email VARCHAR(255) UNIQUE,          
    provider VARCHAR(255),        
    provider_id VARCHAR(255),        
    name VARCHAR(255),   
    locale VARCHAR(255),                 
    profile_picture TEXT,             
    created_at TIMESTAMP DEFAULT NOW(),   
    last_login TIMESTAMP                  
);

CREATE TABLE shops (
    shop_id BIGSERIAL PRIMARY KEY,
    owner_id UUID NOT NULL,
    title VARCHAR(50) NOT NULL,
    domain VARCHAR(50) UNIQUE NOT NULL,
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

CREATE TABLE shop_images (
    shop_image_id BIGSERIAL PRIMARY KEY,
    favicon_url TEXT,
    logo_url TEXT,
    banner_url TEXT,
    cover_image_url TEXT,
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

CREATE TABLE product_types (
    product_type_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    shippable BOOLEAN NOT NULL DEFAULT TRUE,
    digital BOOLEAN NOT NULL DEFAULT FALSE,
    shop_id BIGINT NOT NULL,
    UNIQUE (title, shop_id),
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE products(
    product_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_type_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (title, shop_id),
    CONSTRAINT fk_product_type FOREIGN KEY (product_type_id) REFERENCES product_types(product_type_id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE product_images(
    product_image_id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    alt VARCHAR(255) NOT NULL,
    product_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TYPE product_status AS ENUM('DRAFT','PUBLISHED', 'ARCHIVED');
CREATE TABLE product_variations(
    product_variation_id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    available_quantity BIGINT NOT NULL,
    status product_status NOT NULL DEFAULT 'DRAFT'::product_status,
    seo_description TEXT,
    seo_keywords TEXT[],
    seo_title VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    UNIQUE (slug, shop_id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE 
);

CREATE TYPE attribute_data_type AS ENUM('Text', 'Number', 'Date', 'Option');
CREATE TYPE attribute_unit AS ENUM('KG', 'GB', 'INCH');
CREATE TYPE attribute_applies_to AS ENUM('Product', 'ProductVariation');
CREATE TABLE attributes(
    attribute_id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    data_type attribute_data_type NOT NULL DEFAULT 'Text'::attribute_data_type,
    unit attribute_unit,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    applies_to attribute_applies_to NOT NULL DEFAULT 'Product'::attribute_applies_to,
    shop_id BIGINT NOT NULL,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE attribute_options(
    attribute_option_id BIGSERIAL PRIMARY KEY,
    value VARCHAR(50) NOT NULL,
    shop_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE product_attribute_values(
    product_attribute_value_id BIGSERIAL PRIMARY KEY,
    value VARCHAR(50),
    attribute_option_id BIGINT,
    product_id BIGINT NOT NULL,
    attribute_id BIGINT NOT NULL,
    shop_id BIGINT NOT NULL,
    CONSTRAINT fk_attribute_option FOREIGN KEY (attribute_option_id) REFERENCES attribute_options(attribute_option_id) ON DELETE CASCADE,
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
    CONSTRAINT fk_attribute_option FOREIGN KEY (attribute_option_id) REFERENCES attribute_options(attribute_option_id) ON DELETE CASCADE,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES attributes(attribute_id) ON DELETE CASCADE,
    CONSTRAINT fk_shop FOREIGN KEY (shop_id) REFERENCES shops(shop_id) ON DELETE CASCADE
);

CREATE TABLE shopping_cart(
    shopping_cart_id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE shopping_cart_items(
    shopping_cart_item_id BIGSERIAL PRIMARY KEY,
    quantity BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_variation_id BIGINT NOT NULL,
    shopping_cart_id BIGINT NOT NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_shopping_cart FOREIGN KEY (shopping_cart_id) REFERENCES shopping_cart(shopping_cart_id) ON DELETE CASCADE
);

CREATE TABLE orders(
    order_id BIGSERIAL PRIMARY KEY,
    status VARCHAR(10) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    -- discount_price DECIMAL(10, 2) NOT NULL,
    -- shipping_price DECIMAL(10, 2) NOT NULL,
    -- tax_price DECIMAL(10, 2) NOT NULL,
    -- shipping_address TEXT NOT NULL,
    -- payment_method VARCHAR(10) NOT NULL,
    -- payment_status VARCHAR(10) NOT NULL,
    -- shipping_method VARCHAR(10) NOT NULL,
    -- shipping_status VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE order_items(
    order_item_id BIGSERIAL PRIMARY KEY,
    quantity BIGINT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    product_variation_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    CONSTRAINT fk_product_variation FOREIGN KEY (product_variation_id) REFERENCES product_variations(product_variation_id) ON DELETE CASCADE,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);
