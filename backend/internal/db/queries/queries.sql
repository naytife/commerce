-- name: UpsertUser :one
INSERT INTO users (sub, provider_id, provider, email, name, locale, profile_picture, verified_email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (email)
DO UPDATE SET
    name = COALESCE(EXCLUDED.name, users.name),
    profile_picture = COALESCE(EXCLUDED.profile_picture, users.profile_picture),
    locale = COALESCE(EXCLUDED.locale, users.locale),
    verified_email = COALESCE(EXCLUDED.verified_email, users.verified_email),
    last_login = NOW()
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetUserBySub :one
SELECT * FROM users
WHERE sub = $1;

-- name: GetUserBySubWithShops :one
SELECT 
    users.*,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'shop_id', shops.shop_id,
                'title', shops.title,
                'domain', shops.domain,
                'subdomain', shops.subdomain,
                'status', shops.status,
                'created_at', shops.created_at,
                'updated_at', shops.updated_at
            )
        ) FILTER (WHERE shops.shop_id IS NOT NULL), '[]'::jsonb
    )::jsonb AS shops
FROM users
LEFT JOIN shops ON users.user_id = shops.owner_id
WHERE users.sub = $1
GROUP BY users.user_id;
