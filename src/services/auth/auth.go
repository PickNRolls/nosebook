package auth

import "github.com/google/uuid"

type Auth struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
}
