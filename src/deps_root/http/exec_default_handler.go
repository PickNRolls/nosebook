package roothttp

import (
	"fmt"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/services/auth"
	commandresult "nosebook/src/services/command_result"

	"github.com/gin-gonic/gin"
)

func execDefaultHandler[C any](
	command *C,
	serviceMethod func(*C, *auth.Auth) *commandresult.Result,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		if command != nil {
			if err := ctx.ShouldBindJSON(command); err != nil {
				fmt.Println(err)
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
