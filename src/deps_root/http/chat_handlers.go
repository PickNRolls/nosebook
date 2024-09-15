package roothttp

import (
	presenterchat "nosebook/src/application/presenters/chat"
	presenterdto "nosebook/src/application/presenters/dto"
	presentermessage "nosebook/src/application/presenters/message"
	presenteruser "nosebook/src/application/presenters/user"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addChatHandlers() {
	userPresenter := presenteruser.New(this.db).
		WithTracer(this.tracer)

	messagePresenter := presentermessage.New(this.db, userPresenter).
    WithTracer(this.tracer)
  
	presenter := presenterchat.New(this.db, userPresenter, messagePresenter).
    WithTracer(this.tracer) 

	group := this.authRouter.Group("/chats")

	group.GET("", execDefaultPresenter(presenter.FindByFilter, map[string]presenterOption{
		"interlocutorId": {
			Type: STRING,
		},
		"next": {
			Type: STRING,
		},
		"limit": {
			Type: UINT64,
		},
	}, &presenterchat.FindByFilterInput{}, this.tracer))

	group.GET("/:id", func(ctx *gin.Context) {
		// TODO: apply generic presenter handler
		reqctx := reqcontext.From(ctx)

		output := presenter.FindByFilter(ctx.Request.Context(), &presenterchat.FindByFilterInput{
			Id:    ctx.Param("id"),
			Limit: 1,
		}, reqctx.Auth())
		_, ok := handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		var chat presenterdto.Chat
		if len(output.Data) > 0 {
			chat = output.Data[0]
		}

		reqctx.SetResponseData(chat)
		reqctx.SetResponseOk(true)
	})
}
