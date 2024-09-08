package like

import (
	domainlike "nosebook/src/domain/like"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Repository interface {
	WithPostId(id uuid.UUID) Repository
	WithCommentId(id uuid.UUID) Repository
	WithUserId(id uuid.UUID) Repository
	FindOne() (*domainlike.Like, *Error)

	Save(like *domainlike.Like) *Error
}

type Notifier interface {
	NotifyAbout(like *domainlike.Like) *errors.Error
}

type NotifierRepository interface {
	FindByUserId(id uuid.UUID) Notifier
}
