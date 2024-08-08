package helpers

import (
	"net/http"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
)

func GetAuthOrForbidden(ctx *gin.Context) *auth.Auth {
	a, _ := GetAuthOr(ctx, func() {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
	})

	return a
}
