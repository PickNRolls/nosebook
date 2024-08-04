package commands

import "github.com/google/uuid"

type DenyFriendRequestCommand struct {
	RequesterId uuid.UUID `json:"requesterId"`
}
