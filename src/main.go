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

	"nosebook/src/infra/postgres"
	"nosebook/src/services/user_authentication"
	"nosebook/src/services/user_authentication/commands"
)

func main() {
	postgresPasswordBytes, passwordErr := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if passwordErr != nil {
		log.Fatalln(passwordErr)
	}
	postgresPassword := string(postgresPasswordBytes[:len(postgresPasswordBytes)-1])
	fmt.Printf(postgresPassword)

	db, dbErr := sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:%s@db:5432/%s", postgresPassword, os.Getenv("POSTGRES_DB")))
	if dbErr != nil {
		log.Fatalln(dbErr)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalln(pingErr)
	}

	fmt.Println("Connected to postgres database!")

	userRepository := postgres.NewUserRepository(db)
	userAuthenticationService := services.NewUserAuthenticationService(userRepository)

	router := gin.Default()
	userList := []users.User{}
	router.GET("/", func(ctx *gin.Context) {
		if err := db.Select(&userList, "SELECT * FROM users"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.IndentedJSON(http.StatusOK, userList)
	})

	router.POST("/register", func(ctx *gin.Context) {
		var command commands.RegisterUserCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userAuthenticationService.RegisterUser(&command)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	})

	router.Run("0.0.0.0:8080")
}
