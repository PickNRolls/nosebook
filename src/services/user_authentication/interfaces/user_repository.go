package interfaces

import (
	"nosebook/src/domain/users"
)

type UserRepository interface {
	Create(user *users.User) (*users.User, error)
	FindByNick(nick string) *users.User
}
