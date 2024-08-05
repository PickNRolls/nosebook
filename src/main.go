package main

import (
	"net/http"
	"nosebook/src/domain/users"
	"nosebook/src/handlers"
	"nosebook/src/handlers/friendship"
	"nosebook/src/handlers/posts"

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

	router := gin.Default()

	router.Use(middlewares.NewSessionMiddleware(userAuthenticationService))

	router.GET("/", func(ctx *gin.Context) {
		u := []users.User{}
		if err := db.Select(&u, "SELECT * FROM users"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusOK, u)
	})

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
		group.POST("/publish", posts.NewHandlerPublish(postingService))
		group.POST("/remove", posts.NewHandlerRemove(postingService))
		group.POST("/comment", posts.NewHandlerComment(commentService))
	}

	router.Run("0.0.0.0:8080")
}
