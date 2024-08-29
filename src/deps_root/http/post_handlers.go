package roothttp

import (
	rootpostservice "nosebook/src/deps_root/post_service"
	"nosebook/src/services/posting"
)

func (this *RootHTTP) addPostHandlers() {
	service := rootpostservice.New(this.db)

	group := this.authRouter.Group("/posts")
	group.POST("/publish", execResultHandler(&posting.PublishPostCommand{}, service.Publish))
	group.POST("/remove", execResultHandler(&posting.RemovePostCommand{}, service.Remove))
}
