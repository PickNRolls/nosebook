package roothttp

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"
	commandresult "nosebook/src/lib/command_result"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
)

func execDefaultHandler[C any, T any](
	command *C,
	serviceMethod func(*C, *auth.Auth) (T, *errors.Error),
) func(*gin.Context) {
	return execResultHandler(command, func(c *C, a *auth.Auth) *commandresult.Result {
		return commandresult.FromCommandReturn(serviceMethod(c, a))
	})
}

func execResultHandler[C any](
	command *C,
	serviceMethod func(*C, *auth.Auth) *commandresult.Result,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		if command != nil {
			if err := ctx.ShouldBindJSON(command); err != nil {
				ctx.Error(err)
				ctx.Abort()
				return
			}
		}

		result := serviceMethod(command, reqCtx.Auth())

		reqCtx.SetResponseOk(result.Ok)
		if result.Errors != nil {
			for _, err := range result.Errors {
				ctx.Error(err)
			}
			ctx.Abort()
		}
		reqCtx.SetResponseData(result.Data)
	}
}
