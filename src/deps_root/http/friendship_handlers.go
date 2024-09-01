package roothttp

import (
	"nosebook/src/application/services/friendship"
	rootfriendshipservice "nosebook/src/deps_root/friendship_service"
)

func (this *RootHTTP) addFriendshipHandlers() {
	service := rootfriendshipservice.New(this.db)

	group := this.authRouter.Group("/friendship")
	group.POST("/send-request", execDefaultHandler(&friendship.SendRequestCommand{}, service.SendRequest))
	group.POST("/accept-request", execDefaultHandler(&friendship.AcceptRequestCommand{}, service.AcceptRequest))
	group.POST("/deny-request", execDefaultHandler(&friendship.DenyRequestCommand{}, service.DenyRequest))
	group.POST("/remove-friend", execDefaultHandler(&friendship.RemoveFriendCommand{}, service.RemoveFriend))
}
