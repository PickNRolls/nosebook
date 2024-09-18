package middleware

import (
	userauth "nosebook/src/application/services/user_auth"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func NewSession(service *userauth.Service, tracer trace.Tracer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer ctx.Next()

    _, span := tracer.Start(ctx.Request.Context(), "session_middleware")
    defer span.End()
    
		sessionHeader := ctx.GetHeader("X-Auth-Session-Id")
		if sessionHeader == "" {
			return
		}

		sessionId, err := uuid.Parse(sessionHeader)
		if err != nil {
			return
		}

		user, err := service.TryGetUserBySessionId(userauth.TryGetUserBySessionIdCommand{
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
