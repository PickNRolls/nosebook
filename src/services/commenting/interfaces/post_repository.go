package interfaces

import (
	"nosebook/src/domain/posts"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindById(postId uuid.UUID) *posts.Post
}
