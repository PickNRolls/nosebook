package users

import (
	"errors"
	"nosebook/src/services"
	"nosebook/src/services/user_service/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerGet(userService *services.UserService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.Error(errors.New("No user id"))
			ctx.Abort()
			return
		}
		uuid, err := uuid.Parse(id)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		user, err := userService.GetUser(&commands.GetUserCommand{
			Id: uuid,
		})
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("data", user)
	}
}
