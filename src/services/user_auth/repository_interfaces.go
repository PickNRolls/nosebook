package userauth

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/user"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	FindById(id uuid.UUID) *sessions.Session
	FindByUserId(id uuid.UUID) *sessions.Session
	Create(session *sessions.Session) (*sessions.Session, error)
	Update(session *sessions.Session) (*sessions.Session, error)
	Remove(id uuid.UUID) (*sessions.Session, error)
}

type UserRepository interface {
	Create(user *domainuser.User) (*domainuser.User, error)
	UpdateActivity(userId uuid.UUID, t time.Time) error
	FindByNick(nick string) *domainuser.User
	FindById(id uuid.UUID) *domainuser.User
	FindAll() ([]*domainuser.User, error)
}
