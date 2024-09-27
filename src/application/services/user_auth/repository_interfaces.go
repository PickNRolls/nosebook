package userauth

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/user"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type SessionRepository interface {
	FindById(id uuid.UUID) *sessions.Session
	Create(session *sessions.Session) (*sessions.Session, error)
	Update(session *sessions.Session) (*sessions.Session, error)
	Remove(id uuid.UUID) (*sessions.Session, error)
}

type UserRepository interface {
	Create(user *domainuser.User) (*domainuser.User, error)
	Save(user *domainuser.User) *errors.Error
	FindByNick(nick string) *domainuser.User
	FindById(id uuid.UUID) *domainuser.User
	FindAll() ([]*domainuser.User, error)
}
