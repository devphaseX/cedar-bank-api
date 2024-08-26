package db

import (
	"context"
	"testing"
	"time"

	"github.com/devphasex/cedar-bank-api/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		Fullname:       util.RandomOwner(),
		HashedPassword: "secret",
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)

	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUserByUniqueID(context.Background(), GetUserByUniqueIDParams{
		ID: pgtype.Int8{
			Valid: true,
			Int64: user.ID,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Fullname, user2.Fullname)
	require.Equal(t, user.Email, user2.Email)
	require.WithinDuration(t, user.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
