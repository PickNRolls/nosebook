package middlewares

import (
	"errors"
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewNotAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, ok := helpers.GetUserOr(ctx, func() {
			ctx.Next()
		})

		if ok {
			ctx.Error(errors.New("Only unauthorized users can do it"))
			ctx.Abort()
		}
	}
}
