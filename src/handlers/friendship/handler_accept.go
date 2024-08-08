package friendship

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/friendship/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerAccept(friendshipService *services.FriendshipService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)

		var command commands.AcceptFriendRequestCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request, err := friendshipService.AcceptFriendRequest(&command, &auth.Auth{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, request)
	}
}
