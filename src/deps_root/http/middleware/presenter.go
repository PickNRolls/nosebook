package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	commandresult "nosebook/src/lib/command_result"

	"github.com/gin-gonic/gin"
)

func NewPresenter() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}

		reqctx := reqcontext.From(ctx)

		ctx.JSON(ctx.Writer.Status(), &commandresult.Result{
			Ok:     reqctx.ResponseOk(),
			Errors: reqctx.Errors(),
			Data:   reqctx.ResponseData(),
		})
	}
}
