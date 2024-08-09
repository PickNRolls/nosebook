package handlers

import (
	"net/http"
	"nosebook/src/services"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerLogin(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var command commands.LoginCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authResult, err := userAuthenticationService.Login(&command)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, authResult)
	}
}
