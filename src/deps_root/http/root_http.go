package roothttp

import (
	"fmt"
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/deps_root/http/middleware"
	reqcontext "nosebook/src/deps_root/http/req_context"
	repos "nosebook/src/infra/postgres/repositories"
	"nosebook/src/infra/rabbitmq"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type RootHTTP struct {
	db           *sqlx.DB
	rmqCh        *rabbitmq.Channel
	router       *gin.Engine
	authRouter   *gin.RouterGroup
	unauthRouter *gin.RouterGroup
  traceProvider *trace.TracerProvider
  tracer oteltrace.Tracer
}

func New(db *sqlx.DB, rmqCh *rabbitmq.Channel) *RootHTTP {
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if param.ErrorMessage != "" {
			return fmt.Sprintf("[%s] \"%s %s %d \"%s\" %s %s\"\n",
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.StatusCode,
				param.Request.Header,
				param.Latency,
				param.ErrorMessage,
			)
		}

		return fmt.Sprintf("[%s] \"%s %s %d %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/metrics", func(ctx *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	output := &RootHTTP{
		db:     db,
		router: router,
		rmqCh:  rmqCh,
	}

	userAuthService := userauth.New(repos.NewUserRepository(db), repos.NewSessionRepository(db))
  
	router.Use(middleware.NewRequestMetrics())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(200)
		ctx.Writer.Write([]byte("pong"))
	})
  
  output.enableTracing()

	router.Use(middleware.NewPresenter())
	router.NoRoute(middleware.NewNoRouteHandler())
	router.Use(middleware.NewSession(userAuthService, output.tracer))

	output.unauthRouter = router.Group("/", middleware.NewNotAuth())
	unauthRouter := output.unauthRouter
	output.authRouter = router.Group("/", middleware.NewAuth())
	authRouter := output.authRouter

	unauthRouter.POST("/register", execResultHandler(&userauth.RegisterUserCommand{}, userAuthService.RegisterUser))
	unauthRouter.POST("/login", execResultHandler(&userauth.LoginCommand{}, userAuthService.Login))

	authRouter.POST("/logout", execResultHandler(nil, userAuthService.Logout))
	authRouter.GET("/whoami", func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		user := reqCtx.UserOrForbidden()
		reqCtx.SetResponseOk(true)
		reqCtx.SetResponseData(user)
	})

	output.addLikeHandlers()
	output.addPostHandlers()
	output.addCommentHandlers()
	output.addFriendshipHandlers()
	output.addUserHandlers()
	output.addWebsocketHandlers()
	output.addConversationHandlers()
	output.addChatHandlers()
	output.addMessageHandlers()

	registerMetrics()

	return output
}

func (this *RootHTTP) Run(port string) error {
  err := this.router.Run(port)
  defer this.shutdown()
  return err
}
