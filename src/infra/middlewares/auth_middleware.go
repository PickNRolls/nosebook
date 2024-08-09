package middlewares

import (
	"errors"
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, ok := helpers.GetUserOr(ctx, func() {
			ctx.Error(errors.New("You are not authorized"))
			ctx.Abort()
		})
		if !ok {
			return
		}

		ctx.Next()
	}
}
