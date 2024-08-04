package main

import (
	"fmt"
	"log"
	"net/http"
	"nosebook/src/domain/users"
	"nosebook/src/handlers"
	"os"

	"nosebook/src/infra/middlewares"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"nosebook/src/infra/postgres"
	"nosebook/src/services/user_authentication"
)

func main() {
	var err error

	postgresPasswordBytes, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		log.Fatalln(err)
	}
	postgresPassword := string(postgresPasswordBytes[:len(postgresPasswordBytes)-1])

	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:%s@db:5432/%s", postgresPassword, os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

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
