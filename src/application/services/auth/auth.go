package auth

import (
	"nosebook/src/domain/user"

	"github.com/google/uuid"
)

type Auth struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
}

func From(user *domainuser.User, sessionId uuid.UUID) *Auth {
	if user == nil || sessionId == uuid.Nil {
		return nil
	}

	return &Auth{
		UserId:    user.Id,
		SessionId: sessionId,
	}
}
