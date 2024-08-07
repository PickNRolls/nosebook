package interfaces

import (
	"nosebook/src/domain/posts"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindById(id uuid.UUID) *posts.Post
	Create(post *posts.Post) (*posts.Post, error)
	Save(post *posts.Post) (*posts.Post, error)
	Remove(post *posts.Post) (*posts.Post, error)
}
