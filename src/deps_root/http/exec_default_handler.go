package roothttp

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
)

func execDefaultHandler[C any, R any](
	command C,
	serviceMethod func(*C, *auth.Auth) (R, *errors.Error),
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		response, err := serviceMethod(&command, reqCtx.Auth())
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("data", response)
	}
}
