package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	var (
		txResult TransferTxResult
		err      error
	)

	err = s.execTx(ctx, func(*Queries) error {
		var err error
		fromAccount, err := s.GetAccountByID(ctx, arg.FromAccountID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New(
					fmt.Sprintf("account with from_account_id '%d' not found", arg.FromAccountID),
				)
			}

			return err
		}

		if fromAccount.Balance < arg.Amount {
			return errors.New("sender balance not sufficient for this transaction")
		}

		toAccount, err := s.GetAccountByID(ctx, arg.ToAccountID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.New(
					fmt.Sprintf("account with to_account_id '%d' not found", arg.ToAccountID),
				)
			}

			return err
		}

		tpArg := CreateTransferParams{
			FromAccountID: pgtype.Int8{
				Int64: arg.FromAccountID,
				Valid: true,
			},
			ToAccountID: pgtype.Int8{
				Int64: arg.ToAccountID,
				Valid: true,
			},
			Amount: arg.Amount,
		}

		transfer, err := s.CreateTransfer(ctx, tpArg)

		if err != nil {
			return err
		}

		fromEntryBalanceEntryParams := CreateBalanceEntryParams{
			Amount: -arg.Amount,
			AccountID: pgtype.Int8{
				Int64: fromAccount.ID,
				Valid: true,
			},
		}

		fromEntry, err := s.CreateBalanceEntry(ctx, fromEntryBalanceEntryParams)

		if err != nil {
			return err
		}

		toEntryBalanceEntryParams := CreateBalanceEntryParams{
			Amount: arg.Amount,
			AccountID: pgtype.Int8{
				Int64: toAccount.ID,
				Valid: true,
			},
		}

		toEntry, err := s.CreateBalanceEntry(ctx, toEntryBalanceEntryParams)

		if err != nil {
			return err
		}

		fromAccountBalance := UpdateBalanceParams{
			Balance: fromAccount.Balance - arg.Amount,
			ID:      fromAccount.ID,
		}

		fromAccount, err = s.UpdateBalance(ctx, fromAccountBalance)

		if err != nil {
			return err
		}

		toAccountBalance := UpdateBalanceParams{
			Balance: toAccount.Balance + arg.Amount,
			ID:      toAccount.ID,
		}

		toAccount, err = s.UpdateBalance(ctx, toAccountBalance)

		if err != nil {
			return err
		}

		txResult.FromAccount = fromAccount
		txResult.ToAccount = toAccount
		txResult.FromEntry = fromEntry
		txResult.ToEntry = toEntry
		txResult.Transfer = transfer

		return err
	})

	if err != nil {
		return nil, err
	}

	return &txResult, nil
}
