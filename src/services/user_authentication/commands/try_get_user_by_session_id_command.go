package commands

import "github.com/google/uuid"

type TryGetUserBySessionIdCommand struct {
	SessionId uuid.UUID
}
