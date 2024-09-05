package roothttp

import (
	"nosebook/src/application/services/socket"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addWebsocketHandlers() {
	hub := socket.NewHub()

	group := this.authRouter.Group("/chat")

	group.GET("/", func(ctx *gin.Context) {
		client := socket.NewClient(hub)
		client.Run(ctx)
	})
}
