package roothttp

import (
	"nosebook/src/application/services/auth"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"
	commandresult "nosebook/src/lib/command_result"

	"github.com/gin-gonic/gin"
)

type execHandlerOption interface {
	ShouldAvoidBinding() bool
}

type execHandlerAvoidBinding struct{}

func (this *execHandlerAvoidBinding) ShouldAvoidBinding() bool { return true }

func execDefaultHandler[C any, T any](
	serviceMethod func(C, *auth.Auth) (T, *errors.Error),
	opts ...execHandlerOption,
) func(*gin.Context) {
	return execResultHandler(func(c C, a *auth.Auth) *commandresult.Result {
		return commandresult.FromCommandReturn(serviceMethod(c, a))
	}, opts...)
}

func execResultHandler[C any](
	serviceMethod func(C, *auth.Auth) *commandresult.Result,
	opts ...execHandlerOption,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		var command C

		shouldAvoidBinding := false
		for _, opt := range opts {
			if opt.ShouldAvoidBinding() {
				shouldAvoidBinding = true
			}
		}

		if !shouldAvoidBinding {
			if err := ctx.ShouldBindJSON(&command); err != nil {
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
