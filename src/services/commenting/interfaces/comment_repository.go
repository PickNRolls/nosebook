package interfaces

import (
	"nosebook/src/domain/comments"

	"github.com/google/uuid"
)

type CommentRepository interface {
	FindById(id uuid.UUID) *comments.Comment
	Create(comment *comments.Comment) (*comments.Comment, error)
	Save(comment *comments.Comment) (*comments.Comment, error)
	Remove(id uuid.UUID) (*comments.Comment, error)
}
