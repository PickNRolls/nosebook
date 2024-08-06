package comments

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerRemove(commentService *services.CommentService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrBadRequest(ctx)

		var command commands.RemoveCommentCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		comment, err := commentService.Remove(&command, &auth.Auth{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, comment)
	}
}
