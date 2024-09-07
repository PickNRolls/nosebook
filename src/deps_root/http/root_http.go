package roothttp

import (
	"fmt"
	"nosebook/src/application/services/socket"
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/deps_root/http/middleware"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/infra/postgres/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RootHTTP struct {
	db           *sqlx.DB
	router       *gin.Engine
	authRouter   *gin.RouterGroup
	unauthRouter *gin.RouterGroup
	hub          *socket.Hub
}

func New(db *sqlx.DB, hub *socket.Hub) *RootHTTP {
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

	output := &RootHTTP{
		db:     db,
		router: router,
		hub:    hub,
	}

	userAuthService := userauth.New(repos.NewUserRepository(db), repos.NewSessionRepository(db))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(200)
		ctx.Writer.Write([]byte("pong"))
	})

	router.Use(middleware.NewPresenter())
	router.NoRoute(middleware.NewNoRouteHandler())
	router.Use(middleware.NewSession(userAuthService))

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

	return output
}

func (this *RootHTTP) Run(port string) error {
	return this.router.Run(port)
}
