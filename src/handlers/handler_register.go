package handlers

import (
	"net/http"
	"nosebook/src/services/user_authentication"
	"nosebook/src/services/user_authentication/commands"

	"github.com/gin-gonic/gin"
)

func NewHandlerRegister(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var command commands.RegisterUserCommand
		if err := ctx.ShouldBindJSON(&command); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userAuthenticationService.RegisterUser(&command)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)

		session, err := userAuthenticationService.RegenerateSession(&commands.RegenerateSessionCommand{
			UserId: user.ID,
		})
		if err != nil {
			ctx.Error(err)
		} else {
			ctx.SetCookie("nosebook_session", session.Value.String(), 60*60, "/", "localhost", true, true)
		}
	}
}
