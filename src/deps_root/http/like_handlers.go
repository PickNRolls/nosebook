package roothttp

import (
	rootlikeservice "nosebook/src/deps_root/like_service"
	"nosebook/src/services/like"
)

func (this *RootHTTP) addLikeHandlers() {
	service := rootlikeservice.New(this.db)

	group := this.authRouter.Group("/like")
	group.POST("/post", execResultHandler(&like.LikePostCommand{}, service.LikePost))
	group.POST("/comment", execResultHandler(&like.LikeCommentCommand{}, service.LikeComment))
}
