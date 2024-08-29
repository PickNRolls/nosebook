package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	userauth "nosebook/src/services/user_auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewSession(service *userauth.Service) func(ctx *gin.Context) {
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

		user, err := service.TryGetUserBySessionId(&userauth.TryGetUserBySessionIdCommand{
			SessionId: sessionId,
		})
		if err != nil {
			return
		}

		reqCtx := reqcontext.From(ctx)
		reqCtx.SetUser(user)
		reqCtx.SetSessionId(sessionId)

		if err := service.MarkSessionActive(sessionId); err != nil {
			ctx.Error(err)
		}
	}
}
