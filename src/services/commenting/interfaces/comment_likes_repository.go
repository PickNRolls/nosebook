package interfaces

import (
	"nosebook/src/domain/comments"

	"github.com/google/uuid"
)

type CommentLikesRepository interface {
	Find(commentId uuid.UUID, authorId uuid.UUID) *comments.CommentLike
	Create(like *comments.CommentLike) (*comments.CommentLike, error)
	Remove(like *comments.CommentLike) (*comments.CommentLike, error)
}
