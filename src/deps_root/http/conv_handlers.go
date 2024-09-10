package roothttp

import (
	"nosebook/src/application/services/conversation"
	rootconvservice "nosebook/src/deps_root/conv_service"
)

func (this *RootHTTP) addConversationHandlers() {
	service := rootconvservice.New(this.db, this.rmqCh)

	group := this.authRouter.Group("/conversations")

	group.POST("/send-message", execDefaultHandler(&conversation.SendMessageCommand{}, service.SendMessage))
}
