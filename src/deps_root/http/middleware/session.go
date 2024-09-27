package middleware

import (
	userauth "nosebook/src/application/services/user_auth"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func NewSession(service *userauth.Service, tracer trace.Tracer) func(ctx *gin.Context) {
	return func(ginctx *gin.Context) {
		defer ginctx.Next()

		ctx, span := tracer.Start(ginctx.Request.Context(), "session_middleware")
		defer span.End()

		sessionHeader := ginctx.GetHeader("X-Auth-Session-Id")
		if sessionHeader == "" {
			return
		}

		sessionId, err := uuid.Parse(sessionHeader)
		if err != nil {
			return
		}

		_, span = tracer.Start(ctx, "try_get_user_by_session_id")
		user, err := service.TryGetUserBySessionId(userauth.TryGetUserBySessionIdCommand{
			SessionId: sessionId,
		})
		span.End()
		if err != nil {
			return
		}

		reqctx := reqcontext.From(ginctx)
		reqctx.SetUser(user)
		reqctx.SetSessionId(sessionId)

		_, span = tracer.Start(ctx, "mark_session_active")
		if err := service.MarkSessionActive(ctx, sessionId); err != nil {
			ginctx.Error(err)
		}
		span.End()
	}
}
