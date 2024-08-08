package interfaces

import (
	"nosebook/src/domain/users"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *users.User) (*users.User, error)
	UpdateActivity(userId uuid.UUID, t time.Time) error
	FindByNick(nick string) *users.User
	FindById(id uuid.UUID) *users.User
	FindAll() ([]*users.User, error)
}
