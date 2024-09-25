package roothttp

import (
	"fmt"
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/deps_root/http/middleware"
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
	db            *sqlx.DB
	rmqConn       *rabbitmq.Connection
	router        *gin.Engine
	authRouter    *gin.RouterGroup
	unauthRouter  *gin.RouterGroup
	traceProvider *trace.TracerProvider
	tracer        oteltrace.Tracer
	shutdowns     []ShutdownFn
}

func New(db *sqlx.DB, rmqConn *rabbitmq.Connection) *RootHTTP {
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
		db:      db,
		router:  router,
		rmqConn: rmqConn,
	}

	router.Use(middleware.NewRequestMetrics())
  router.Use(middleware.NewDbMetrics(db))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(200)
		ctx.Writer.Write([]byte("pong"))
	})

	output.enableTracing()

	sessionRepository := repos.NewSessionRepository(db)
	go sessionRepository.Run()
  output.shutdowns = append(output.shutdowns, sessionRepository.OnDispose)
	userAuthService := userauth.New(repos.NewUserRepository(db), sessionRepository, output.tracer)

	router.Use(middleware.NewPresenter())
	router.NoRoute(middleware.NewNoRouteHandler())
	router.Use(middleware.NewSession(userAuthService, output.tracer))

	output.unauthRouter = router.Group("/", middleware.NewNotAuth())
	output.authRouter = router.Group("/", middleware.NewAuth())

	output.addAuthHandlers(userAuthService)
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
