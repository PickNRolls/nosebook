package roothttp

import (
	rootpostservice "nosebook/src/deps_root/post_service"
	"nosebook/src/services/posting"
)

func (this *RootHTTP) addPostHandlers() {
	service := rootpostservice.New(this.db)

	group := this.authRouter.Group("/posts")
	group.POST("/publish", execDefaultHandler(posting.PublishPostCommand{}, service.Publish))
	group.POST("/remove", execDefaultHandler(posting.RemovePostCommand{}, service.Remove))
}
