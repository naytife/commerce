-- name: UpsertUser :one
INSERT INTO users (auth0_sub, email, name, profile_picture_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (auth0_sub)
DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name, profile_picture_url = EXCLUDED.profile_picture_url
RETURNING user_id, auth0_sub, email, name, profile_picture_url;

-- name: GetUser :one
SELECT * FROM users
WHERE auth0_sub = $1;