package roothttp

import (
	"nosebook/src/application/services/socket"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addWebsocketHandlers() {
	group := this.authRouter.Group("/ws")

	group.GET("/", func(ctx *gin.Context) {
		client := socket.NewClient(this.hub)
		client.Run(ctx)
	})
}
