package helpers

import (
	"net/http"
	"nosebook/src/domain/users"
	"nosebook/src/infra/errors"

	"github.com/gin-gonic/gin"
)

func GetUserOrForbidden(ctx *gin.Context) *users.User {
	user, _ := GetUserOr(ctx, func() {
		ctx.Status(http.StatusForbidden)
		ctx.Error(errors.NewNotAuthenticatedError())
		ctx.Abort()
	})

	return user
}
