package roothttp

import (
	presentercomment "nosebook/src/application/presenters/comment"
	rootcommentpresenter "nosebook/src/deps_root/comment_presenter"
	rootcommentservice "nosebook/src/deps_root/comment_service"
	"nosebook/src/deps_root/http/exec"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addCommentHandlers() {
	service := rootcommentservice.New(this.db)
	presenter := rootcommentpresenter.New(this.db, this.tracer)

	group := this.authRouter.Group("/comments")

	group.POST("/publish-on-post", exec.Command(service.PublishOnPost, exec.WithUuidMapper))
	group.POST("/remove", exec.Command(service.Remove, exec.WithUuidMapper))

	group.GET("", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(ctx.Request.Context(), presentercomment.FindByFilterInput{
			PostId: ctx.Query("postId"),
			Next:   ctx.Query("next"),
			Prev:   ctx.Query("prev"),
			Last:   ctx.Query("last") == "true",
			Limit:  5,
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
