package interfaces

import (
	"nosebook/src/domain/friendship"

	"github.com/google/uuid"
)

type UserFriendsRepository interface {
	FindByBoth(requesterId uuid.UUID, responderId uuid.UUID) *friendship.FriendRequest
	Create(friendRequest *friendship.FriendRequest) (*friendship.FriendRequest, error)
	Update(friendRequest *friendship.FriendRequest) (*friendship.FriendRequest, error)
}
