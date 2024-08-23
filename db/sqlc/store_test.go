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

	// fmt.Printf("before: acct1 %v acct2 %v\n", account1.Balance, account2.Balance)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		ctx := context.Background()
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

	exist := map[int]bool{}
	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotNil(t, result)
		// fmt.Printf("after: acct1 %v acct2 %v\n", result.FromAccount.Balance, result.ToAccount.Balance)

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
		require.NotContains(t, exist, k)
		exist[k] = true

	}
	updatedAccount1, err := store.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := store.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-float64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+float64(n)*amount, updatedAccount2.Balance)

}

func TestFundTransferDeadlock(t *testing.T) {
	store := testQueries

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	var amount float64 = 10

	errs := make(chan error, n)

	// fmt.Printf("before: acct1 %v acct2 %v\n", account1.Balance, account2.Balance)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		ctx := context.Background()
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err

		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errs)

	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

	}
	updatedAccount1, err := store.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := store.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
