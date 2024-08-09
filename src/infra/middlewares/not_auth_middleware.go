package middlewares

import (
	"net/http"
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewNotAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, ok := helpers.GetUserOr(ctx, func() {
			ctx.Next()
		})

		if ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Only unauthorized users can do it"})
		}
	}
}
