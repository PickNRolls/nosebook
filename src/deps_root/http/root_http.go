package roothttp

import (
	"nosebook/src/deps_root/http/middleware"
	"nosebook/src/handlers"
	"nosebook/src/infra/postgres/repositories"
	"nosebook/src/services"

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

	userAuthenticationService := services.NewUserAuthenticationService(repos.NewUserRepository(db), repos.NewSessionRepository(db))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(200)
		ctx.Writer.Write([]byte("pong"))
	})

	router.Use(middleware.NewPresenter())
	router.NoRoute(middleware.NewNoRouteHandler())
	router.Use(middleware.NewSession(userAuthenticationService))

	output.unauthRouter = router.Group("/", middleware.NewNotAuth())
	unauthRouter := output.unauthRouter
	output.authRouter = router.Group("/", middleware.NewAuth())
	authRouter := output.authRouter

	unauthRouter.POST("/register", handlers.NewHandlerRegister(userAuthenticationService))
	unauthRouter.POST("/login", handlers.NewHandlerLogin(userAuthenticationService))

	authRouter.GET("/whoami", handlers.NewHandlerWhoAmI())
	authRouter.POST("/logout", handlers.NewHandlerLogout(userAuthenticationService))

	output.addLikeHandlers()
	output.addPostHandlers()
	output.addCommentHandlers()
	output.addFriendshipHandlers()

	return output
}

func (this *RootHTTP) Run(port string) error {
	return this.router.Run(port)
}
