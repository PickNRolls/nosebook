package auth

import (
	"nosebook/src/domain/users"

	"github.com/google/uuid"
)

type Auth struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
}

func From(user *users.User, sessionId uuid.UUID) *Auth {
	if user == nil || sessionId == uuid.Nil {
		return nil
	}

	return &Auth{
		UserId:    user.ID,
		SessionId: sessionId,
	}
}
