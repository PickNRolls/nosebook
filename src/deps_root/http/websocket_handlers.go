package roothttp

import (
	"log"
	"net/http"
	"net/http/httputil"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var InProgressWsConnectionsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "ws_connections_in_progress",
		Help: "Number of current websocket connections proxying into Notification service",
	},
)

func (this *RootHTTP) addWebsocketHandlers() {
	group := this.authRouter.Group("/ws")

	group.GET("", func(ctx *gin.Context) {
		auth := reqcontext.From(ctx).Auth()

		log.Printf("Proxying websocket connection to notification service for user(id:%v)\n", auth.UserId)

		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "notification:8081"
			req.Header["X-Auth-User-Id"] = []string{auth.UserId.String()}
		}}

    InProgressWsConnectionsGauge.Inc();
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
    InProgressWsConnectionsGauge.Dec();
	})
}
