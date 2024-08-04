package middlewares

import (
	"nosebook/src/services/user_authentication"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSessionMiddleware(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer ctx.Next()

		sessionCookie, err := ctx.Cookie("nosebook_session")
		if err != nil {
			return
		}

		sessionId, err := uuid.Parse(sessionCookie)
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
		ctx.SetCookie("nosebook_session", sessionId.String(), 60*60, "/", "localhost", true, true)
	}
}
