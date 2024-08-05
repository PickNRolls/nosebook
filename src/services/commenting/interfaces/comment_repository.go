package interfaces

import (
	"nosebook/src/domain/comments"

	"github.com/google/uuid"
)

type CommentRepository interface {
	FindById(id uuid.UUID) *comments.Comment
	CreateForPost(postId uuid.UUID, comment *comments.Comment) (*comments.Comment, error)
	Remove(id uuid.UUID) (*comments.Comment, error)
}
