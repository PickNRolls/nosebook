package main

import (
	"nosebook/src/handlers"
	"nosebook/src/handlers/comments"
	"nosebook/src/handlers/friendship"
	"nosebook/src/handlers/posts"
	"nosebook/src/handlers/users"
	"nosebook/src/presenters"

	"nosebook/src/infra/middlewares"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"nosebook/src/infra/postgres"
	"nosebook/src/infra/postgres/repositories"
	"nosebook/src/services"
)

func main() {
	db := postgres.Connect()

	userAuthenticationService := services.NewUserAuthenticationService(repos.NewUserRepository(db), repos.NewSessionRepository(db))
	friendshipService := services.NewFriendshipService(repos.NewUserFriendsRepository(db))
	postingService := services.NewPostingService(repos.NewPostsRepository(db))
	commentService := services.NewCommentService(repos.NewCommentRepository(db))
	userService := services.NewUserService(repos.NewUserRepository(db))

	postPresenter := presenters.
		NewPostPresenter().
		WithPostingService(postingService).
		WithPostRepository(repos.NewPostPresenterRepository(db)).
		WithCommentService(commentService)

	router := gin.Default()

	router.Use(middlewares.NewPresenterMiddleware())
	router.Use(middlewares.NewErrorHandlerMiddleware())
	router.Use(middlewares.NewSessionMiddleware(userAuthenticationService))

	unauthRouter := router.Group("/", middlewares.NewNotAuthMiddleware())
	unauthRouter.POST("/register", handlers.NewHandlerRegister(userAuthenticationService))
	unauthRouter.POST("/login", handlers.NewHandlerLogin(userAuthenticationService))

	authRouter := router.Group("/")
	authRouter.Use(middlewares.NewAuthMiddleware())

	authRouter.GET("/whoami", handlers.NewHandlerWhoAmI())
	authRouter.POST("/logout", handlers.NewHandlerLogout(userAuthenticationService))

	{
		group := authRouter.Group("/friendship")
		group.POST("/add", friendship.NewHandlerAdd(friendshipService))
		group.POST("/accept", friendship.NewHandlerAccept(friendshipService))
		group.POST("/deny", friendship.NewHandlerDeny(friendshipService))
		group.POST("/remove", friendship.NewHandlerRemove(friendshipService))
	}

	{
		group := authRouter.Group("/posts")
		group.GET("/", posts.NewHandlerFind(postPresenter))

		group.POST("/publish", posts.NewHandlerPublish(postPresenter))
		group.POST("/remove", posts.NewHandlerRemove(postPresenter))
		group.POST("/like", posts.NewHandlerLike(postPresenter))
	}

	{
		group := authRouter.Group("/comments")
		group.GET("/", comments.NewHandlerFind(commentService))

		group.POST("/publish-on-post", comments.NewHandlerPublishOnPost(commentService))
		group.POST("/remove", comments.NewHandlerRemove(commentService))
		group.POST("/like", comments.NewHandlerLike(commentService))
	}

	{
		group := authRouter.Group("/users")
		group.GET("/", users.NewHandlerGetAll(userService))
		group.GET("/:id", users.NewHandlerGet(userService))
	}

	router.Run("0.0.0.0:8080")
}
