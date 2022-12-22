package jwt

import (
	"errors"
	"time"
)

var (
	ErrExpired = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        int       `json:""`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id int, duration time.Duration) *Payload {
	return &Payload{
		ID: id,
		IssueAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpired
	}

	return nil
}