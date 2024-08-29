package roothttp

import (
	"nosebook/src/handlers/friendship"
	"nosebook/src/infra/postgres/repositories"
	"nosebook/src/services"
)

func (this *RootHTTP) addFriendshipHandlers() {
	friendshipService := services.NewFriendshipService(repos.NewUserFriendsRepository(this.db))

	group := this.authRouter.Group("/friendship")
	group.POST("/add", friendship.NewHandlerAdd(friendshipService))
	group.POST("/accept", friendship.NewHandlerAccept(friendshipService))
	group.POST("/deny", friendship.NewHandlerDeny(friendshipService))
	group.POST("/remove", friendship.NewHandlerRemove(friendshipService))
}
