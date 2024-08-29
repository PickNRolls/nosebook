package roothttp

import (
	"nosebook/src/deps_root/http/middleware"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/infra/postgres/repositories"
	userauth "nosebook/src/services/user_auth"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RootHTTP struct {
	db           *sqlx.DB
	router       *gin.Engine
	authRouter   *gin.RouterGroup
	unauthRouter *gin.RouterGroup
}

func New(db *sqlx.DB) *RootHTTP {
	router := gin.Default()
	output := &RootHTTP{
		db:     db,
		router: router,
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

	unauthRouter.POST("/register", execDefaultHandler(&userauth.RegisterUserCommand{}, userAuthService.RegisterUser))
	unauthRouter.POST("/login", execDefaultHandler(&userauth.LoginCommand{}, userAuthService.Login))

	authRouter.POST("/logout", execDefaultHandler(nil, userAuthService.Logout))
	authRouter.GET("/whoami", func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		user := reqCtx.UserOrForbidden()
		reqCtx.SetResponseOk(true)
		reqCtx.SetResponseData(user)
	})

	output.addLikeHandlers()
	output.addPostHandlers()
	output.addCommentHandlers()
	// output.addFriendshipHandlers()

	return output
}

func (this *RootHTTP) Run(port string) error {
	return this.router.Run(port)
}
