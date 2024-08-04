package main

import (
	"fmt"
	"log"
	"net/http"
	"nosebook/src/domain/users"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"nosebook/src/infra/middlewares"

	"nosebook/src/infra/postgres"
	"nosebook/src/services/user_authentication"
	"nosebook/src/services/user_authentication/commands"
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

	router.POST("/register", func(ctx *gin.Context) {
		var command commands.RegisterUserCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userAuthenticationService.RegisterUser(&command)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)

		session, err := userAuthenticationService.RegenerateSession(&commands.RegenerateSessionCommand{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
		} else {
			ctx.SetCookie("nosebook_session", session.Value.String(), 60*60, "/", "localhost", true, true)
		}
	})

	router.Run("0.0.0.0:8080")
}
