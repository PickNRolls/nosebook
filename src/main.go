package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
}

func main() {
	users := []User{}

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

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		db.Select(&users, "SELECT * FROM users")
		ctx.IndentedJSON(http.StatusOK, users)
	})

	router.Run("0.0.0.0:8080")
}
