package helpers

import (
	"nosebook/src/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserOr(ctx *gin.Context, fn func()) (*users.User, bool) {
	userAny, ok := ctx.Get("user")
	if !ok {
		fn()
		return nil, false
	}

	user, ok := userAny.(*users.User)
	if !ok {
		fn()
		return nil, false
	}

	return user, true
}
