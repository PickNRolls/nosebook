package main

import (
	"fmt"
	"log"
	"net/http"

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

	db, err := sqlx.Connect("pgx", "postgres://postgres:123@database:5432/nosebook")
	if err != nil {
		log.Fatalln(err)
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
