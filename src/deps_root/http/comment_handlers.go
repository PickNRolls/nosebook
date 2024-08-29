package roothttp

import (
	rootcommentservice "nosebook/src/deps_root/comment_service"
	"nosebook/src/services/commenting"
)

func (this *RootHTTP) addCommentHandlers() {
	commentService := rootcommentservice.New(this.db)

	group := this.authRouter.Group("/comments")

	group.POST("/publish-on-post", execDefaultHandler(&commenting.PublishPostCommentCommand{}, commentService.PublishOnPost))
	group.POST("/remove", execDefaultHandler(&commenting.RemoveCommentCommand{}, commentService.Remove))
}
