package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paster       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, errors.New(fmt.Sprintf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize))
	}

	maker := &PasetoMaker{
		symmetricKey: []byte(symmetricKey),
		paster:       paseto.NewV2(),
	}

	return maker, nil
}

func (p *PasetoMaker) CreateToken(userId int64, email string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userId, email, duration)
	if err != nil {
		return "", nil, err
	}

	var tokenStr string
	tokenStr, err = p.paster.Encrypt(p.symmetricKey, payload, nil)

	return tokenStr, payload, nil
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var (
		payload Payload
		err     error
	)

	if err = p.paster.Decrypt(token, p.symmetricKey, &payload, nil); err != nil {
		return nil, ErrInvalidToken
	}

	if err = payload.Valid(); err != nil {
		return nil, err
	}

	return &payload, nil
}
