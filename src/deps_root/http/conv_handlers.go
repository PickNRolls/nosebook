package roothttp

import (
	presentermessage "nosebook/src/application/presenters/message"
	rootconvservice "nosebook/src/deps_root/conv_service"
)

func (this *RootHTTP) addConversationHandlers(presenter *presentermessage.Presenter) {
	service := rootconvservice.New(this.db, this.rmqConn, presenter, this.tracer)

	group := this.authRouter.Group("/conversations")

	group.POST("/send-message", execCommand(service.SendMessage, this))
}
