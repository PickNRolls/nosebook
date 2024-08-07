package users

import (
	"net/http"
	"nosebook/src/services"
	"nosebook/src/services/user_service/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerGet(userService *services.UserService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id"})
			return
		}
		uuid, err := uuid.Parse(id)

		user, err := userService.GetUser(&commands.GetUserCommand{
			Id: uuid,
		})
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
