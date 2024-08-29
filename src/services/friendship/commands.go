package friendship

import "github.com/google/uuid"

type AcceptRequestCommand struct {
	RequesterId uuid.UUID `json:"requesterId"`
}

type DenyRequestCommand struct {
	RequesterId uuid.UUID `json:"requesterId"`
}

type RemoveFriendCommand struct {
	FriendId uuid.UUID `json:"friendId"`
}

type SendRequestCommand struct {
	ResponderId uuid.UUID `json:"responderId"`
	Message     string    `json:"message"`
}
