package handlers

import (
	"nosebook/src/services"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerRegister(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var command commands.RegisterUserCommand
		err := ctx.ShouldBindJSON(&command)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		authResult, err := userAuthenticationService.RegisterUser(&command)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Set("data", authResult)
	}
}
