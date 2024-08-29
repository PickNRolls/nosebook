package posting

import "github.com/google/uuid"

type EditPostCommand struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

type PublishPostCommand struct {
	Message string    `json:"message"`
	OwnerId uuid.UUID `json:"ownerId"`
}

type RemovePostCommand struct {
	Id uuid.UUID `json:"id"`
}
