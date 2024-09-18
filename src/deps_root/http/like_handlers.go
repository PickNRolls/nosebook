package roothttp

import (
	rootlikeservice "nosebook/src/deps_root/like_service"
)

func (this *RootHTTP) addLikeHandlers() {
	service := rootlikeservice.New(this.db, this.rmqCh)

	group := this.authRouter.Group("/like")
	group.POST("/post", execDefaultHandler(service.LikePost))
	group.POST("/comment", execDefaultHandler(service.LikeComment))
}
