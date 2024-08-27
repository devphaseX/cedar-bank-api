package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken          = errors.New("token has expired")
	ErrInvalidToken          = errors.New("token not valid")
	ErrInvalidOrExpiredToken = errors.New("token expired or invalid")
	ErrUnverifiableToken     = errors.New("token is unverifiable")
)

type Payload struct {
	ID     uuid.UUID `json:"id"`
	UserId int64     `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func NewPayload(userId int64, email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:     tokenID,
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt.Time) {
		return ErrExpiredToken
	}

	return nil
}
