package commands

import "github.com/google/uuid"

type EditPostCommand struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}
