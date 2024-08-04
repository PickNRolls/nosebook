package main

import (
	"net/http"
	"nosebook/src/domain/users"
	"nosebook/src/handlers"

	"nosebook/src/infra/middlewares"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"nosebook/src/infra/postgres"
	"nosebook/src/services/user_authentication"
)

func main() {
	db := postgres.Connect()

	userRepository := postgres.NewUserRepository(db)
	sessionRepository := postgres.NewSessionRepository(db)
	userAuthenticationService := services.NewUserAuthenticationService(userRepository, sessionRepository)

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

	router.Run("0.0.0.0:8080")
}
