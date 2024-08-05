package posts

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerPublish(postingService *services.PostingService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user, ok := helpers.GetUserOrBadRequest(ctx)
		if !ok {
			return
		}

		var command commands.PublishPostCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request, err := postingService.Publish(&command, &auth.Auth{
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
