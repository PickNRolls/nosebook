package user

import (
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	"nosebook/src/lib/image"

	"github.com/google/uuid"
)

type AvatarStorage interface {
	Upload(image *image.Image, userId uuid.UUID) (string, *errors.Error)
}

type UserRepository interface {
	FindById(id uuid.UUID) *domainuser.User
	Save(user *domainuser.User) *errors.Error
}
