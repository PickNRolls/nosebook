package commands

import "github.com/google/uuid"

type RemoveFriendCommand struct {
	FriendId uuid.UUID `json:"friendId"`
}
