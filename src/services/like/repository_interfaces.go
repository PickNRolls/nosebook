package like

import (
	domainlike "nosebook/src/domain/like"

	"github.com/google/uuid"
)

type Repository interface {
	WithPostId(id uuid.UUID) Repository
	WithCommentId(id uuid.UUID) Repository
	WithUserId(id uuid.UUID) Repository
	FindOne() (*domainlike.Like, *Error)

	Save(like *domainlike.Like) *Error
}
