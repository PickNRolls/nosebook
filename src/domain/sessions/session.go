package sessions

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Value          uuid.UUID `db:"session"`
	UserId         uuid.UUID `db:"user_id"`
	CreatedAt      time.Time `db:"created_at"`
	LastActivityAt time.Time `db:"last_activity_at"`
}

func NewSession(userId uuid.UUID) *Session {
	now := time.Now()

	return &Session{
		Value:          uuid.New(),
		UserId:         userId,
		CreatedAt:      now,
		LastActivityAt: now,
	}
}
