package commands

import "github.com/google/uuid"

type RemovePostCommand struct {
	Id uuid.UUID `json:"id"`
}
