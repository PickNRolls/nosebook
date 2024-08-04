package commands

import "github.com/google/uuid"

type RegenerateSessionCommand struct {
	UserId uuid.UUID
}
