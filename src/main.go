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
	"nosebook/src/services"
)

func main() {
	db := postgres.Connect()

	userAuthenticationService := services.NewUserAuthenticationService(postgres.NewUserRepository(db), postgres.NewSessionRepository(db))
	friendshipService := services.NewFriendshipService(postgres.NewUserFriendsRepository(db))
	postingService := services.NewPostingService(postgres.NewPostsRepository(db))
	commentService := services.NewCommentService(postgres.NewCommentRepository(db))
	userService := services.NewUserService(postgres.NewUserRepository(db))

	postPresenter := presenters.NewPostPresenter().WithPostingService(postingService).WithPostRepository(postgres.NewPostPresenterRepository(db))

	router := gin.Default()

	router.Use(middlewares.NewSessionMiddleware(userAuthenticationService))

	router.POST("/register", handlers.NewHandlerRegister(userAuthenticationService))

	authRouter := router.Group("/")
	authRouter.Use(middlewares.NewAuthMiddleware())

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

		group.POST("/publish", posts.NewHandlerPublish(postingService))
		group.POST("/remove", posts.NewHandlerRemove(postingService))
		group.POST("/like", posts.NewHandlerLike(postingService))
	}

	{
		group := authRouter.Group("/comments")
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
