package users

import (
	"nosebook/src/services"

	"github.com/gin-gonic/gin"
)

func NewHandlerGetAll(userService *services.UserService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		users, err := userService.GetAllUsers()
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("data", users)
	}
}
