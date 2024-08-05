package commands

import "github.com/google/uuid"

type PublishPostCommand struct {
	Message string    `json:"message"`
	OwnerId uuid.UUID `json:"ownerId"`
}
