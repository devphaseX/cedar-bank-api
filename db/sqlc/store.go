package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/devphasex/cedar-bank-api/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error)
}

type PgStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PgStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *PgStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	db := New(tx)

	err = fn(db)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var ErrFundNotSufficient = util.NewCustomError("ErrFundNotSufficient", "insufficient funds for transfer")
var ErrUnableUpdateAccount = util.NewCustomError("ErrUnableUpdateAccount", "failed to update both accounts")
var ErrAccountNotFound = util.NewCustomError("ErrAccountNotFound", "account not found")

func (s *PgStore) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	var txResult TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		fromAccount, err := s.GetAccountByIDForUpdate(ctx, arg.FromAccountID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Join(errors.New(
					fmt.Sprintf("account with '%d' not found", arg.FromAccountID),
				), ErrAccountNotFound)
			}

			return err
		}

		// Check if the sender has sufficient funds
		senderBalance := fromAccount.Balance

		if senderBalance < arg.Amount {
			return ErrFundNotSufficient
		}

		_, err = s.GetAccountByIDForUpdate(ctx, arg.ToAccountID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Join(errors.New(
					fmt.Sprintf("account with '%d' not found", arg.ToAccountID),
				), ErrAccountNotFound)
			}

			return err
		}

		// Create transfer
		transfer, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
			ToAccountID:   pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		txResult.Transfer = transfer

		// Create entries
		txResult.FromEntry, err = q.CreateBalanceEntry(ctx, CreateBalanceEntryParams{
			AccountID: pgtype.Int8{Int64: arg.FromAccountID, Valid: true},
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		txResult.ToEntry, err = q.CreateBalanceEntry(ctx, CreateBalanceEntryParams{
			AccountID: pgtype.Int8{Int64: arg.ToAccountID, Valid: true},
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		updateTransferAccountBal := UpdateTransferAccountBalanceParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}
		// Update both account balances in a single query

		result, err := q.UpdateTransferAccountBalance(ctx, updateTransferAccountBal)
		if err != nil {
			return err
		}

		if result.RowsAffected() != 2 {
			return errors.New("failed to update both accounts")
		}

		// Fetch updated accounts
		txResult.FromAccount, err = q.GetAccountByID(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		txResult.ToAccount, err = q.GetAccountByID(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &txResult, nil
}
