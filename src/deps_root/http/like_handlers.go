package roothttp

import (
	"nosebook/src/application/services/like"
	rootlikeservice "nosebook/src/deps_root/like_service"
)

func (this *RootHTTP) addLikeHandlers() {
	service := rootlikeservice.New(this.db, this.hub)

	group := this.authRouter.Group("/like")
	group.POST("/post", execDefaultHandler(&like.LikePostCommand{}, service.LikePost))
	group.POST("/comment", execDefaultHandler(&like.LikeCommentCommand{}, service.LikeComment))
}
