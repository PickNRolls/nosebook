package helpers

import (
	"net/http"
	"nosebook/src/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserOrBadRequest(ctx *gin.Context) *users.User {
	user, _ := GetUserOr(ctx, func() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No user"})
	})

	return user
}
