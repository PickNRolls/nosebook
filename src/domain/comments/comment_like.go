package comments

import (
	"github.com/google/uuid"
)

type CommentLike struct {
	CommentId uuid.UUID `db:"comment_id"`
	AuthorId  uuid.UUID `db:"author_id"`
}

func NewCommentLike(commentId uuid.UUID, authorId uuid.UUID) *CommentLike {
	return &CommentLike{
		AuthorId:  authorId,
		CommentId: commentId,
	}
}
