package db

import (
	"context"
	"testing"
	"time"

	"github.com/devphasex/cedar-bank-api/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	require.NotEmpty(t, user)
	arg := CreateAccountParams{
		OwnerID:  user.ID,
		Balance:  float64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerID, account.OwnerID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountByID(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerID, account2.OwnerID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateBalanceParams{
		Balance: float64(util.RandomMoney()),
		ID:      account1.ID,
	}

	for account1.Balance == arg.Balance {
		arg.Balance = float64(util.RandomMoney())
	}

	account2, err := testQueries.UpdateBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.OwnerID, account2.OwnerID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccountByID(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	accounts := make([]Account, 0, 10)
	l := cap(accounts)
	for i := 0; i < l; i++ {
		account := createRandomAccount(t)
		accounts = append(accounts, account)
	}

	arg := GetAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	queriedAccounts, err := testQueries.GetAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, queriedAccounts, 5)

	for _, account := range queriedAccounts {
		require.NotEmpty(t, account)
	}
}
