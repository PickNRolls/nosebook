package posts

import (
	"github.com/google/uuid"
)

type PostLike struct {
	PostId   uuid.UUID `db:"post_id"`
	AuthorId uuid.UUID `db:"author_id"`
}

func NewPostLike(postId uuid.UUID, authorId uuid.UUID) *PostLike {
	return &PostLike{
		AuthorId: authorId,
		PostId:   postId,
	}
}
