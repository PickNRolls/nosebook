package interfaces

import (
	"nosebook/src/domain/users"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *users.User) (*users.User, error)
	FindByNick(nick string) *users.User
	FindById(id uuid.UUID) *users.User
	FindAll() ([]*users.User, error)
}
