package token

import "time"

type Maker interface {
	CreateToken(userId int64, email string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
