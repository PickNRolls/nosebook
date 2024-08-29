package commenting

import (
	"nosebook/src/domain/comment"
	"nosebook/src/domain/post"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Repository interface {
	FindById(id uuid.UUID, includeRemoved bool) *domaincomment.Comment
	Save(comment *domaincomment.Comment) *errors.Error
}

type PostRepository interface {
	FindById(id uuid.UUID) *domainpost.Post
}
