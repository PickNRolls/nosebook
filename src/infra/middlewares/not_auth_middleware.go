package middlewares

import (
	"nosebook/src/infra/errors"
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewNotAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, ok := helpers.GetUserOr(ctx, func() {
			ctx.Next()
		})

		if ok {
			ctx.Error(errors.NewAuthorizedError())
			ctx.Abort()
		}
	}
}
