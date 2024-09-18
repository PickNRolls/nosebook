package roothttp

import (
	presenterfriendship "nosebook/src/application/presenters/friendship"
	rootfriendshippresenter "nosebook/src/deps_root/friendship_presenter"
	rootfriendshipservice "nosebook/src/deps_root/friendship_service"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addFriendshipHandlers() {
	service := rootfriendshipservice.New(this.db)
	presenter := rootfriendshippresenter.New(this.db)

	group := this.authRouter.Group("/friendship")
	group.POST("/send-request", execDefaultHandler(service.SendRequest))
	group.POST("/accept-request", execDefaultHandler(service.AcceptRequest))
	group.POST("/deny-request", execDefaultHandler(service.DenyRequest))
	group.POST("/remove-friend", execDefaultHandler(service.RemoveFriend))

	group.GET("", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)
		userId := ctx.Query("userId")
		accepted := reqctx.QueryNullableBool("accepted")
		viewed := reqctx.QueryNullableBool("viewed")
		limit, ok := reqctx.QueryNullableUint64("limit")
		if !ok {
			return
		}

		_, onlyIncoming := ctx.GetQuery("onlyIncoming")
		_, onlyOutcoming := ctx.GetQuery("onlyOutcoming")
		_, onlyOnline := ctx.GetQuery("onlyOnline")

		output := presenter.FindByFilter(presenterfriendship.FindByFilterInput{
			UserId:   userId,
			Accepted: accepted,
			Viewed:   viewed,

			OnlyIncoming:  onlyIncoming,
			OnlyOutcoming: onlyOutcoming,
			OnlyOnline:    onlyOnline,

			Limit: limit.Value,
			Next:  ctx.Query("next"),
			Prev:  ctx.Query("prev"),
		}, reqctx.Auth())
		_, ok = handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})

	group.GET("/relation-between", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)
		sourceUserId := ctx.Query("sourceUserId")
		targetUserIds := ctx.QueryArray("targetUserIds")

		out, ok := handle(presenter.DescribeRelation(presenterfriendship.DescribeRelationInput{
			SourceUserId:  sourceUserId,
			TargetUserIds: targetUserIds,
		}, reqctx.Auth()))(reqctx)

		if !ok {
			return
		}

		reqctx.SetResponseData(out)
		reqctx.SetResponseOk(true)
	})
}
