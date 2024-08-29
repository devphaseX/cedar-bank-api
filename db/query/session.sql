-- name: CreateSession :one
INSERT INTO sessions(
id,
owner_id,
user_agent,
refresh_token,
client_ip,
is_blocked,
expired_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSessionList :many
SELECT * FROM sessions
where (sqlc.narg('id')::uuid IS NULL OR sqlc.narg('id')::uuid = sessions.id)
AND (sqlc.narg('owner_id')::bigint IS NULL OR sqlc.narg('owner_id')::bigint  = sessions.owner_id);


-- name: GetSessionByUniqueID :one
SELECT * FROM sessions
where (sqlc.narg('id')::uuid IS NULL OR sqlc.narg('id')::uuid = sessions.id)
AND (sqlc.narg('owner_id')::bigint IS NULL OR sqlc.narg('owner_id')::bigint  = sessions.owner_id)
AND (sqlc.narg('refresh_token')::text IS NULL OR sqlc.narg('refresh_token')::text  = sessions.refresh_token)

LIMIT 1;
