package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New(fmt.Sprintf("Invalid key size: must be at least %d characters", minSecretKeySize))
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userId int64, email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userId, email, duration)

	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var payload *Payload
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is unverifiable") {
			return nil, ErrUnverifiableToken
		}

		return nil, ErrInvalidOrExpiredToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
