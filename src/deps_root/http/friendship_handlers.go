package roothttp

import (
	presenterfriendship "nosebook/src/application/presenters/friendship"
	"nosebook/src/application/services/friendship"
	rootfriendshippresenter "nosebook/src/deps_root/friendship_presenter"
	rootfriendshipservice "nosebook/src/deps_root/friendship_service"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addFriendshipHandlers() {
	service := rootfriendshipservice.New(this.db)
	presenter := rootfriendshippresenter.New(this.db)

	group := this.authRouter.Group("/friendship")
	group.POST("/send-request", execDefaultHandler(&friendship.SendRequestCommand{}, service.SendRequest))
	group.POST("/accept-request", execDefaultHandler(&friendship.AcceptRequestCommand{}, service.AcceptRequest))
	group.POST("/deny-request", execDefaultHandler(&friendship.DenyRequestCommand{}, service.DenyRequest))
	group.POST("/remove-friend", execDefaultHandler(&friendship.RemoveFriendCommand{}, service.RemoveFriend))

	group.GET("/", func(ctx *gin.Context) {
		userId := ctx.Query("userId")
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(&presenterfriendship.FindByFilterInput{
			UserId: userId,
		}, reqctx.Auth())
		_, ok := handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})
}
