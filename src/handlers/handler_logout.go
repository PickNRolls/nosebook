package handlers

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/services"

	"github.com/gin-gonic/gin"
)

func NewHandlerLogout(userAuthenticationService *services.UserAuthenticationService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		auth := helpers.GetAuthOrForbidden(ctx)

		session, err := userAuthenticationService.Logout(auth)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, session)
	}
}
