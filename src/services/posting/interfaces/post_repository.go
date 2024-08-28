package interfaces

import (
	"nosebook/src/domain/posts"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindById(postId uuid.UUID) *posts.Post
	Save(post *posts.Post) *errors.Error
}
