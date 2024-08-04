package commands

import "github.com/google/uuid"

type SendFriendRequestCommand struct {
	ResponderId uuid.UUID `json:"responderId"`
	Message     string    `json:"message"`
}
