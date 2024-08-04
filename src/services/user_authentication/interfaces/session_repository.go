package interfaces

import (
	"nosebook/src/domain/sessions"

	"github.com/google/uuid"
)

type SessionRepository interface {
	FindById(id uuid.UUID) *sessions.Session
	FindByUserId(id uuid.UUID) *sessions.Session
	Create(session *sessions.Session) (*sessions.Session, error)
	Update(session *sessions.Session) (*sessions.Session, error)
}
