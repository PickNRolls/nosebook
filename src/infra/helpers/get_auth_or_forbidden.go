package helpers

import (
	"nosebook/src/infra/errors"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
)

func GetAuthOrForbidden(ctx *gin.Context) *auth.Auth {
	a, _ := GetAuthOr(ctx, func() {
		ctx.Error(errors.NewNotAuthorizedError())
		ctx.Abort()
	})

	return a
}
