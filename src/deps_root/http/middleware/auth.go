package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func NewAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := reqcontext.From(ctx).UserOrForbidden()
		if user == nil {
			return
		}

		ctx.Next()
	}
}
