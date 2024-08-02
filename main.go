package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func main() {
	fmt.Println("Server started ⚡")

	users := []User{
		{"1", "Ilya", "Blinkov"},
		{"2", "Ilya", "Olegblinkov"},
	}

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, users)
	})

	router.Run("localhost:8080")
}
