package middlewares

import (
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)
		if user == nil {
			return
		}

		ctx.Next()
	}
}
