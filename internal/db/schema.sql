CREATE TABLE shops (
    id UUID PRIMARY KEY,
    -- owner_id UUID NOT NULL,
    title VARCHAR(50) NOT NULL,
    default_domain VARCHAR(50) UNIQUE NOT NULL,
    favicon_url TEXT,
    currency_code VARCHAR(3) NOT NULL,
    -- status VARCHAR(10) NOT NULL,
    about TEXT,
    -- seo_description TEXT,
    -- seo_keywords TEXT,
    -- seo_title VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
