-- name: UpsertUser :one
INSERT INTO users (provider_id, provider, email, name, locale, profile_picture)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (email)
DO UPDATE SET
    name = COALESCE(EXCLUDED.name, users.name),
    profile_picture = COALESCE(EXCLUDED.profile_picture, users.profile_picture),
    locale = COALESCE(EXCLUDED.locale, users.locale)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;