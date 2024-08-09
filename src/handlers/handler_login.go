package handlers

import (
	"nosebook/src/services"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerLogin(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var command commands.LoginCommand
		err := ctx.ShouldBindJSON(&command)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		authResult, err := userAuthenticationService.Login(&command)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("data", authResult)
	}
}
