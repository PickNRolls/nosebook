package middlewares

import (
	"nosebook/src/services"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSessionMiddleware(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer ctx.Next()

		sessionHeader := ctx.GetHeader("X-Auth-Session-Id")
		if sessionHeader == "" {
			return
		}

		sessionId, err := uuid.Parse(sessionHeader)
		if err != nil {
			return
		}

		user, err := userAuthenticationService.TryGetUserBySessionId(&commands.TryGetUserBySessionIdCommand{
			SessionId: sessionId,
		})
		if err != nil {
			return
		}

		ctx.Set("user", user)

		if err := userAuthenticationService.MarkSessionActive(sessionId); err != nil {
			ctx.Error(err)
		}
	}
}
