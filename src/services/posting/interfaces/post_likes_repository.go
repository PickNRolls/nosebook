package interfaces

import (
	"nosebook/src/domain/posts"

	"github.com/google/uuid"
)

type PostLikesRepository interface {
	Find(postId uuid.UUID, authorId uuid.UUID) *posts.PostLike
	Create(like *posts.PostLike) (*posts.PostLike, error)
	Remove(like *posts.PostLike) (*posts.PostLike, error)
}
