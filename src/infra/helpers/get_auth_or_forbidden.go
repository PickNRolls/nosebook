package helpers

import (
	"errors"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
)

func GetAuthOrForbidden(ctx *gin.Context) *auth.Auth {
	a, _ := GetAuthOr(ctx, func() {
		ctx.Error(errors.New("You are not authorized"))
		ctx.Abort()
	})

	return a
}
