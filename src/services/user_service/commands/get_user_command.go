package commands

import "github.com/google/uuid"

type GetUserCommand struct {
	Id uuid.UUID `json:"id"`
}
