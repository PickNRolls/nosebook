package interfaces

import (
	"nosebook/src/domain/comments"

	"github.com/google/uuid"
)

type CommentRepository interface {
	FindById(id uuid.UUID, includeRemoved bool) *comments.Comment
	Create(comment *comments.Comment) (*comments.Comment, error)
	Save(comment *comments.Comment) (*comments.Comment, error)
}
