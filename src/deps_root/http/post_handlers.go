package roothttp

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	rootpostservice "nosebook/src/deps_root/post_service"
	presenterpost "nosebook/src/presenters/post"
	"nosebook/src/services/posting"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addPostHandlers() {
	service := rootpostservice.New(this.db)
	presenter := presenterpost.New(this.db)

	group := this.authRouter.Group("/posts")
	group.POST("/publish", execResultHandler(&posting.PublishPostCommand{}, service.Publish))
	group.POST("/remove", execResultHandler(&posting.RemovePostCommand{}, service.Remove))

	group.GET("/", func(ctx *gin.Context) {
		authorId := ctx.Query("authorId")
		ownerId := ctx.Query("ownerId")
		cursor := ctx.Query("cursor")
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(&presenterpost.FindByFilterInput{
			AuthorId: authorId,
			OwnerId:  ownerId,
			Cursor:   cursor,
		}, reqctx.Auth())

		if output.Err != nil {
			ctx.Error(output.Err)
			ctx.Abort()
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})
}
