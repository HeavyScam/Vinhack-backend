package token

import (
	"time"
)

type Maker interface {
	GenerateToken(email string, gender int, userID string, duration time.Duration) (string, *Payload, error)

	//validate token
	VerifyToken(token string) (*Payload, error)
}
