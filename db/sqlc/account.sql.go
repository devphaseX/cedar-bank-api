// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts(owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING id, owner, balance, currency, created_at
`

type CreateAccountParams struct {
	Owner    string  `json:"owner"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAccount, id)
	return err
}

const getAccountByID = `-- name: GetAccountByID :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAccountByID(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountByID, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountByIDForUpdate = `-- name: GetAccountByIDForUpdate :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = $1
LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetAccountByIDForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountByIDForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccounts = `-- name: GetAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE ($3::int[] IS NULL OR id = ANY($3::int[]))
  AND ($1::int IS NULL OR balance < $1)
OFFSET $2
LIMIT $4
`

type GetAccountsParams struct {
	Balance pgtype.Int4 `json:"balance"`
	Offset  int64       `json:"offset"`
	Column3 []int32     `json:"column_3"`
	Limit   int64       `json:"limit"`
}

func (q *Queries) GetAccounts(ctx context.Context, arg GetAccountsParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAccounts,
		arg.Balance,
		arg.Offset,
		arg.Column3,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBalance = `-- name: UpdateBalance :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING id, owner, balance, currency, created_at
`

type UpdateBalanceParams struct {
	Balance float64 `json:"balance"`
	ID      int64   `json:"id"`
}

func (q *Queries) UpdateBalance(ctx context.Context, arg UpdateBalanceParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateBalance, arg.Balance, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateTransferAccountBalance = `-- name: UpdateTransferAccountBalance :execresult
UPDATE accounts
SET balance =
CASE
    WHEN id = $1 THEN balance - $2
    WHEN id = $3 THEN balance + $2
END
WHERE id IN ($1, $3)
`

type UpdateTransferAccountBalanceParams struct {
	FromAccountID int64   `json:"from_account_id"`
	Amount        float64 `json:"amount"`
	ToAccountID   int64   `json:"to_account_id"`
}

func (q *Queries) UpdateTransferAccountBalance(ctx context.Context, arg UpdateTransferAccountBalanceParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, updateTransferAccountBalance, arg.FromAccountID, arg.Amount, arg.ToAccountID)
}
