package roothttp

import (
	rootconvservice "nosebook/src/deps_root/conv_service"
)

func (this *RootHTTP) addConversationHandlers() {
	service := rootconvservice.New(this.db, this.rmqConn)

	group := this.authRouter.Group("/conversations")

	group.POST("/send-message", execDefaultHandler(service.SendMessage))
}
