package commands

import "github.com/google/uuid"

type CreateSessionCommand struct {
	UserId uuid.UUID
}
