package roothttp

import (
	presentermessage "nosebook/src/application/presenters/message"
	presenteruser "nosebook/src/application/presenters/user"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addMessageHandlers() {
	userPresenter := presenteruser.New(this.db)
	presenter := presentermessage.New(this.db, userPresenter)

	group := this.authRouter.Group("/messages")

	group.GET("", func(ctx *gin.Context) {
		// TODO: create generic declarative presenter handler
		reqctx := reqcontext.From(ctx)

		chatId := ctx.Query("chatId")
		next := ctx.Query("next")
		prev := ctx.Query("prev")
		limit, ok := reqctx.QueryNullableUint64("limit")
		if !ok {
			return
		}

		output := presenter.FindByFilter(&presentermessage.FindByFilterInput{
			ChatId: chatId,
			Next:   next,
			Prev:   prev,
			Limit:  limit.Value,
		}, reqctx.Auth())
		_, ok = handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})
}
