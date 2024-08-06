package users

import (
	"net/http"
	"nosebook/src/services"

	"github.com/gin-gonic/gin"
)

func NewHandlerGetAll(userService *services.UserService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		users, err := userService.GetAllUsers()
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, users)
	}
}
