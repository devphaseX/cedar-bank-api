-- name: CreateAccount :one
INSERT INTO accounts(owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts
WHERE ($3::int[] IS NULL OR id = ANY($3::int[]))
  AND (sqlc.narg('balance')::int IS NULL OR balance < sqlc.narg('balance'))
OFFSET sqlc.arg('offset')
LIMIT sqlc.arg('limit');

-- name: UpdateBalance :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
