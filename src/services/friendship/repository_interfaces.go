package friendship

import (
	"nosebook/src/domain/friendship"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Repository interface {
	RequesterId(id uuid.UUID) Repository
	ResponderId(id uuid.UUID) Repository
	OnlyAccepted() Repository
	OnlyNotAccepted() Repository
	FindOne() *domainfriendship.FriendRequest

	Save(request *domainfriendship.FriendRequest) *errors.Error
}
