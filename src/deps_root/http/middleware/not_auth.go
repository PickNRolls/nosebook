package middleware

import (
	"net/http"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/infra/errors"

	"github.com/gin-gonic/gin"
)

func NewNotAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := reqcontext.From(ctx).User()
		if user == nil {
			ctx.Next()
			return
		}

		ctx.Status(http.StatusForbidden)
		ctx.Error(errors.NewAuthenticatedError())
		ctx.Abort()
	}
}
