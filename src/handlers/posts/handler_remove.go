//go:build exclude

package posts

import (
	"nosebook/src/infra/helpers"
	// "nosebook/src/presenters"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerRemove(postPresenter *presenters.PostPresenter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)

		var command commands.RemovePostCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		post, error := postPresenter.Remove(&command, &auth.Auth{
			UserId: user.ID,
		})
		if error != nil {
			ctx.Error(error)
			ctx.Abort()
			return
		}

		ctx.Set("data", post)
	}
}
