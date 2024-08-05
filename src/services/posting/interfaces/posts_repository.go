package interfaces

import "nosebook/src/domain/posts"

// "nosebook/src/domain/friendship"
//
// "github.com/google/uuid"

type PostsRepository interface {
	Create(post *posts.Post) (*posts.Post, error)
}
