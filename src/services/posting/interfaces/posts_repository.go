package interfaces

import (
	"nosebook/src/domain/posts"

	"github.com/google/uuid"
)

type PostsRepository interface {
	FindById(id uuid.UUID) *posts.Post
	Create(post *posts.Post) (*posts.Post, error)
	Remove(post *posts.Post) (*posts.Post, error)
}
