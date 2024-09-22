package roothttp

import (
	"nosebook/src/deps_root/http/exec"
	rootlikeservice "nosebook/src/deps_root/like_service"
)

func (this *RootHTTP) addLikeHandlers() {
	service := rootlikeservice.New(this.db, this.rmqConn)

	group := this.authRouter.Group("/like")
	group.POST("/post", exec.Command(service.LikePost))
	group.POST("/comment", exec.Command(service.LikeComment))
}
