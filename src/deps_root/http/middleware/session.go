package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/services"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSession(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
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

		reqCtx := reqcontext.From(ctx)
		reqCtx.SetUser(user)
		reqCtx.SetSessionId(sessionId)

		if err := userAuthenticationService.MarkSessionActive(sessionId); err != nil {
			ctx.Error(err)
		}
	}
}
