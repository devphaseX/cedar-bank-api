package token

import "time"

type Maker interface {
	CreateToken(userId int64, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
