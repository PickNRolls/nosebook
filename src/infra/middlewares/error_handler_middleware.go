package middlewares

import (
	"github.com/gin-gonic/gin"
	"nosebook/src/errors"
)

func NewErrorHandlerMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		errList := []*errors.Error{}
		for _, ginErr := range ctx.Errors {
			err := errors.From(ginErr.Err)
			errList = append(errList, err)
		}

		ctx.Set("errors", errList)
	}
}
