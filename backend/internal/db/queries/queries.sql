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