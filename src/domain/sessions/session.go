package sessions

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	Value     uuid.UUID `db:"session"`
	UserId    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewSession(userId uuid.UUID) *Session {
	now := time.Now()

	return &Session{
		Value:     uuid.New(),
		UserId:    userId,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
