package interfaces

import (
	"nosebook/src/domain/comments"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type CommentRepository interface {
	FindById(id uuid.UUID, includeRemoved bool) *comments.Comment
	Save(comment *comments.Comment) *errors.Error
}
