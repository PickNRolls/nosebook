package helpers

import (
	"net/http"
	"nosebook/src/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserOrBadRequest(ctx *gin.Context) (*users.User, bool) {
	userAny, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user"})
		return nil, false
	}

	user, ok := userAny.(*users.User)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No user"})
		return nil, false
	}

	return user, true
}
