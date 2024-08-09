package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewErrorHandlerMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		errors := []error{}
		for _, err := range ctx.Errors {
			fmt.Println(err)
			errors = append(errors, err)
		}

		ctx.Set("errors", errors)
	}
}
