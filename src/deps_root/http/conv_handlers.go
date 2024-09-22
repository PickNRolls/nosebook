package roothttp

import (
	rootconvservice "nosebook/src/deps_root/conv_service"
)

func (this *RootHTTP) addConversationHandlers() {
	service := rootconvservice.New(this.db, this.rmqConn, this.tracer)

	group := this.authRouter.Group("/conversations")

	group.POST("/send-message", execCommand(service.SendMessage, this))
}
