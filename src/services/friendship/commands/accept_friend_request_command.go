package commands

import "github.com/google/uuid"

type AcceptFriendRequestCommand struct {
	RequesterId uuid.UUID `json:"requesterId"`
}
