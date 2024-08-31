package roothttp

import (
	rootcommentpresenter "nosebook/src/deps_root/comment_presenter"
	rootcommentservice "nosebook/src/deps_root/comment_service"
	reqcontext "nosebook/src/deps_root/http/req_context"
	presentercomment "nosebook/src/presenters/comment"
	"nosebook/src/services/commenting"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addCommentHandlers() {
	service := rootcommentservice.New(this.db)
	presenter := rootcommentpresenter.New(this.db)

	group := this.authRouter.Group("/comments")

	group.POST("/publish-on-post", execResultHandler(&commenting.PublishPostCommentCommand{}, service.PublishOnPost))
	group.POST("/remove", execResultHandler(&commenting.RemoveCommentCommand{}, service.Remove))

	group.GET("/", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(&presentercomment.FindByFilterInput{
			PostId: ctx.Query("postId"),
			Next:   ctx.Query("next"),
			Prev:   ctx.Query("prev"),
			Last:   ctx.Query("last") == "true",
			Limit:  5,
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
