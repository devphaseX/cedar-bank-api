// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64              `json:"id"`
	OwnerID   int64              `json:"owner_id"`
	Balance   float64            `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Entry struct {
	ID        int64              `json:"id"`
	AccountID pgtype.Int8        `json:"account_id"`
	Amount    float64            `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Session struct {
	ID           pgtype.UUID        `json:"id"`
	OwnerID      int64              `json:"owner_id"`
	UserAgent    string             `json:"user_agent"`
	RefreshToken string             `json:"refresh_token"`
	ClientIp     pgtype.Text        `json:"client_ip"`
	IsBlocked    pgtype.Bool        `json:"is_blocked"`
	ExpiredAt    pgtype.Timestamptz `json:"expired_at"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

type Transfer struct {
	ID            int64       `json:"id"`
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	// amount must be positive
	Amount    float64            `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type User struct {
	ID                int64              `json:"id"`
	Username          string             `json:"username"`
	Email             string             `json:"email"`
	Fullname          string             `json:"fullname"`
	HashedPassword    string             `json:"hashed_password"`
	PasswordSalt      string             `json:"password_salt"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}
