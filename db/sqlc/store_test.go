package db

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFundTransfer(t *testing.T) {
	store := testQueries

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	var amount float64 = 10

	errs := make(chan error, n)
	results := make(chan *TransferTxResult, n)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errs)
	close(results)

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotNil(t, result)

		// Perform more detailed checks on the result
		if result != nil {
			require.NotEmpty(t, result.Transfer)
			require.Equal(t, account1.ID, result.Transfer.FromAccountID.Int64)
			require.Equal(t, account2.ID, result.Transfer.ToAccountID.Int64)
			require.Equal(t, amount, result.Transfer.Amount)
			require.NotZero(t, result.Transfer.ID)
			require.NotZero(t, result.Transfer.CreatedAt)

			require.NotEmpty(t, result.FromEntry)
			require.Equal(t, account1.ID, result.FromEntry.AccountID.Int64)
			require.Equal(t, -amount, result.FromEntry.Amount)
			require.NotZero(t, result.FromEntry.ID)
			require.NotZero(t, result.FromEntry.CreatedAt)

			require.NotEmpty(t, result.ToEntry)
			require.Equal(t, account2.ID, result.ToEntry.AccountID.Int64)
			require.Equal(t, amount, result.ToEntry.Amount)
			require.NotZero(t, result.ToEntry.ID)
			require.NotZero(t, result.ToEntry.CreatedAt)
			fromAccount := result.FromAccount
			require.NotEmpty(t, fromAccount)
			require.Equal(t, account1.ID, fromAccount.ID)

			toAccount := result.ToAccount
			require.NotEmpty(t, toAccount)
			require.Equal(t, account2.ID, toAccount.ID)

			diff1 := account1.Balance - fromAccount.Balance
			diff2 := toAccount.Balance - account2.Balance
			require.Equal(t, diff1, diff2)
			require.True(t, diff1 > 0)

			k := int(diff1 / amount)
			require.True(t, k >= 1 && k <= n)
		}
	}
}
