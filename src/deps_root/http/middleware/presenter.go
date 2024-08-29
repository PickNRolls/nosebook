package middleware

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"

	"github.com/gin-gonic/gin"
)

type responseDTO struct {
	Errors []*errors.Error `json:"errors"`
	Data   any             `json:"data"`
}

func NewPresenter() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		reqContext := reqcontext.From(ctx)

		responseDTO := &responseDTO{
			Data:   reqContext.ResponseData(),
			Errors: reqContext.Errors(),
		}

		ctx.JSON(ctx.Writer.Status(), responseDTO)
	}
}
