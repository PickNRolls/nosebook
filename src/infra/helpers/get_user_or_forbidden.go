package helpers

import (
	"errors"
	"nosebook/src/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserOrForbidden(ctx *gin.Context) *users.User {
	user, _ := GetUserOr(ctx, func() {
		ctx.Error(errors.New("You are not authorized"))
		ctx.Abort()
	})

	return user
}
