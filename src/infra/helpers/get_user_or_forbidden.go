package helpers

import (
	"nosebook/src/domain/users"
	"nosebook/src/infra/errors"

	"github.com/gin-gonic/gin"
)

func GetUserOrForbidden(ctx *gin.Context) *users.User {
	user, _ := GetUserOr(ctx, func() {
		ctx.Error(errors.NewNotAuthorizedError())
		ctx.Abort()
	})

	return user
}
