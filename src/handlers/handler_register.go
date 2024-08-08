package handlers

import (
	"net/http"
	"nosebook/src/domain/sessions"
	"nosebook/src/domain/users"
	"nosebook/src/services"
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

		session, err := userAuthenticationService.CreateSession(&commands.CreateSessionCommand{
			UserId: user.ID,
		})

		if err != nil {
			ctx.Error(err)
		}

		result := &struct {
			User    *users.User       `json:"user"`
			Session *sessions.Session `json:"session"`
		}{
			User:    user,
			Session: session,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
