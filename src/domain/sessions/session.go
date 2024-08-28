package sessions

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionId uuid.UUID `json:"sessionId" db:"session_id"`
	UserId    uuid.UUID `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

func NewSession(userId uuid.UUID) *Session {
	now := time.Now()

	return &Session{
		SessionId: uuid.New(),
		UserId:    userId,
		CreatedAt: now,
		ExpiresAt: now.Add(48 * time.Hour),
	}
}

func (s *Session) Refresh() error {
	now := time.Now()

	if s.ExpiresAt.Before(now) {
		return errors.New("Can't refresh session, it is expired.")
	}

	s.ExpiresAt = now.Add(48 * time.Hour)
	return nil
}
