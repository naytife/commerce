-- name: UpsertUser :one
INSERT INTO users ( email, name, profile_picture)
VALUES ($1, $2, $3)
ON CONFLICT (email)
DO UPDATE SET name = EXCLUDED.name, profile_picture = EXCLUDED.profile_picture
RETURNING user_id, email, name, profile_picture;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;