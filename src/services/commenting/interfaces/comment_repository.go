package interfaces

import (
	"nosebook/src/domain/comments"
	"nosebook/src/generics"
	"nosebook/src/services/commenting/structs"

	"github.com/google/uuid"
)

type CommentRepository interface {
	FindById(id uuid.UUID, includeRemoved bool) *comments.Comment
	FindByFilter(filter structs.QueryFilter, limitSize *uint) *generics.SingleQueryResult[*comments.Comment]
	Create(comment *comments.Comment) (*comments.Comment, error)
	Save(comment *comments.Comment) (*comments.Comment, error)
}
