-- name: CreateTransfer :one
INSERT INTO transfer(from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING  *;
