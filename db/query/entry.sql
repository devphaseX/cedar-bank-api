-- name: CreateBalanceEntry :one
INSERT INTO entries (account_id, amount)
VALUES($1, $2)
RETURNING *;


-- name: GetBalanceEntry :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;

-- name: GetAccountBalanceEntries :many
SELECT * FROM entries
WHERE account_id = $1
LIMIT 1;
