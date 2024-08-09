package handlers

import (
	"nosebook/src/infra/helpers"

	"github.com/gin-gonic/gin"
)

func NewHandlerWhoAmI() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)
		ctx.Set("data", user)
	}
}
