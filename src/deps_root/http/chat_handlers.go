package roothttp

import (
	presenterchat "nosebook/src/application/presenters/chat"
	presentermessage "nosebook/src/application/presenters/message"
	presenteruser "nosebook/src/application/presenters/user"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addChatHandlers() {
	userPresenter := presenteruser.New(this.db)
	messagePresenter := presentermessage.New(this.db, userPresenter)
	presenter := presenterchat.New(this.db, userPresenter, messagePresenter)

	group := this.authRouter.Group("/chats")

	group.GET("", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)

		next := ctx.Query("next")
		limit, ok := reqctx.QueryNullableUint64("limit")
		if !ok {
			return
		}

		output := presenter.FindByFilter(&presenterchat.FindByFilterInput{
			Next:  next,
			Limit: limit.Value,
		}, reqctx.Auth())
		_, ok = handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	})
}
