package interfaces

import (
	"nosebook/src/domain/posts"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindById(postId uuid.UUID) *posts.Post
	Save(post *posts.Post) *errors.Error
	// FindByFilter(filter structs.QueryFilter) *generics.SingleQueryResult[*posts.Post]
	//
	// Create(post *posts.Post) (*posts.Post, error)
	// Save(post *posts.Post) (*posts.Post, error)
	// Remove(post *posts.Post) (*posts.Post, error)
}
