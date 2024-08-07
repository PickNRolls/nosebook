package posts

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/presenters"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerLike(postPresenter *presenters.PostPresenter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrBadRequest(ctx)

		var command commands.LikePostCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post, err := postPresenter.Like(&command, &auth.Auth{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, post)
	}
}
