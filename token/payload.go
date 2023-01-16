package token

import (
	"errors"
	// "fmt"
	"time"
	//import uuid
	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

//payload is a struct that contains the data that will be stored in the token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	UserID   string    `json:"user_id"`
	Gender   int       `json:"gender"`
	IssueAt  time.Time `json:"issue_at"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewPayload(email string, gender int, userID string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:       tokenID,
		Email:    email,
		UserID:   userID,
		Gender:   gender,
		IssueAt:  time.Now(),
		ExpireAt: time.Now().Add(duration),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}
