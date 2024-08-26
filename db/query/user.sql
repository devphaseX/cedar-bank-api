-- name: CreateUser :one
INSERT INTO users(username, email, fullname, hashed_password, password_salt)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByUniqueID :one
SELECT * FROM users
WHERE id = sqlc.narg('id')
or email ilike sqlc.narg('email')
or username ilike sqlc.narg('username')
LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
WHERE ($3::int[] IS NULL OR id = ANY($3::int[]))
OFFSET sqlc.arg('offset')
LIMIT sqlc.arg('limit');
