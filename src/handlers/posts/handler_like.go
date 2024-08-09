package posts

import (
	"nosebook/src/infra/helpers"
	"nosebook/src/presenters"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerLike(postPresenter *presenters.PostPresenter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)

		var command commands.LikePostCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		post, err := postPresenter.Like(&command, &auth.Auth{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("data", post)
	}
}
