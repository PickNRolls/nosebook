package posting

import (
	"nosebook/src/domain/post"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Repository interface {
	FindById(id uuid.UUID) *domainpost.Post
	Save(post *domainpost.Post) *errors.Error
}
