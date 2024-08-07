package interfaces

import (
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/structs"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindById(id uuid.UUID) *posts.Post
	FindByFilter(filter structs.QueryFilter) structs.QueryResult

	Create(post *posts.Post) (*posts.Post, error)
	Save(post *posts.Post) (*posts.Post, error)
	Remove(post *posts.Post) (*posts.Post, error)
}
