package middlewares

import (
	"net/http"
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, ok := helpers.GetUserOr(ctx, func() {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
		})
		if !ok {
			return
		}

		ctx.Next()
	}
}
