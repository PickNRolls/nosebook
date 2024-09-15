package roothttp

import (
	presenterpost "nosebook/src/application/presenters/post"
	"nosebook/src/application/services/posting"
	reqcontext "nosebook/src/deps_root/http/req_context"
	rootpostpresenter "nosebook/src/deps_root/post_presenter"
	rootpostservice "nosebook/src/deps_root/post_service"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addPostHandlers() {
	service := rootpostservice.New(this.db)
	presenter := rootpostpresenter.New(this.db, this.tracer)

	group := this.authRouter.Group("/posts")
	group.POST("/publish", execResultHandler(&posting.PublishPostCommand{}, service.Publish))
	group.POST("/remove", execResultHandler(&posting.RemovePostCommand{}, service.Remove))

	group.GET("", func(ctx *gin.Context) {
		authorId := ctx.Query("authorId")
		ownerId := ctx.Query("ownerId")
		cursor := ctx.Query("cursor")
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(ctx.Request.Context(), &presenterpost.FindByFilterInput{
			AuthorId: authorId,
			OwnerId:  ownerId,
			Cursor:   cursor,
		}, reqctx.Auth())
		_, ok := handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})

	group.GET("/:id", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)
		id := ctx.Param("id")

		reqctx.SetResponseData(presenter.FindById(ctx.Request.Context(), id, reqctx.Auth()))
		reqctx.SetResponseOk(true)
	})
}
