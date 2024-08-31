package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	commandresult "nosebook/src/lib/command_result"

	"github.com/gin-gonic/gin"
)

func NewPresenter() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		reqContext := reqcontext.From(ctx)

		ctx.JSON(ctx.Writer.Status(), &commandresult.Result{
			Ok:     reqContext.ResponseOk(),
			Errors: reqContext.Errors(),
			Data:   reqContext.ResponseData(),
		})
	}
}
