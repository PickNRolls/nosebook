package helpers

import (
	"net/http"
	"nosebook/src/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserOrForbidden(ctx *gin.Context) *users.User {
	user, _ := GetUserOr(ctx, func() {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized"})
	})

	return user
}
