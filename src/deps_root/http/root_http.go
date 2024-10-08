package roothttp

import (
	"fmt"
	presentermessage "nosebook/src/application/presenters/message"
	presenteruser "nosebook/src/application/presenters/user"
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/deps_root/http/middleware"
	"nosebook/src/domain/sessions"
	domainuser "nosebook/src/domain/user"
	repos "nosebook/src/infra/postgres/repositories"
	userrepo "nosebook/src/infra/postgres/repositories/user_repository"
	"nosebook/src/infra/rabbitmq"
	"nosebook/src/lib/cache"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			return fmt.Sprintf("[%s] \"%s %s %d \"%s\" %s %s\"",
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

	sessionRepository := repos.NewSessionRepository(db).CacheWith(cache.NewLRU[uuid.UUID, *sessions.Session](256))
	output.shutdownRun(sessionRepository)

	userRepository := userrepo.New(db).CacheWith(cache.NewLRU[uuid.UUID, *domainuser.User](256))
	output.shutdownRun(userRepository)

	userAuthService := userauth.New(userRepository, sessionRepository, output.tracer)

	router.Use(middleware.NewPresenter())
	router.NoRoute(middleware.NewNoRouteHandler())
	router.Use(middleware.NewSession(userAuthService, output.tracer))

	output.unauthRouter = router.Group("/", middleware.NewNotAuth())
	output.authRouter = router.Group("/", middleware.NewAuth())

	messagePresenter := presentermessage.New(output.db, presenteruser.New(output.db))
	output.shutdownRun(messagePresenter)

	output.addAuthHandlers(userAuthService)
	output.addLikeHandlers()
	output.addPostHandlers()
	output.addCommentHandlers()
	output.addFriendshipHandlers()
	output.addUserHandlers(userRepository)
	output.addWebsocketHandlers()
	output.addConversationHandlers(messagePresenter)
	output.addChatHandlers(messagePresenter)
	output.addMessageHandlers(messagePresenter)

	registerMetrics()

	return output
}

func (this *RootHTTP) Run(port string) error {
	err := this.router.Run(port)
	defer this.shutdown()
	return err
}
