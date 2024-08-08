package helpers

import (
	"nosebook/src/domain/users"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthOr(ctx *gin.Context, fn func()) (*auth.Auth, bool) {
	userAny, ok := ctx.Get("user")
	if !ok {
		fn()
		return nil, false
	}

	sessionAny, ok := ctx.Get("sessionId")
	if !ok {
		fn()
		return nil, false
	}

	user, ok := userAny.(*users.User)
	if !ok {
		fn()
		return nil, false
	}

	sessionId, ok := sessionAny.(uuid.UUID)
	if !ok {
		fn()
		return nil, false
	}

	return &auth.Auth{
		UserId:    user.ID,
		SessionId: sessionId,
	}, true
}
