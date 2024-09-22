package roothttp

import (
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/deps_root/http/exec"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addAuthHandlers(userAuthService *userauth.Service) {
	this.unauthRouter.POST("/register", exec.Command(userAuthService.RegisterUser))
	this.unauthRouter.POST("/login", exec.Command(userAuthService.Login))

	this.authRouter.POST("/logout", exec.Command(userAuthService.Logout, exec.WithAvoidBinding))
	this.authRouter.GET("/whoami", func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		user := reqCtx.UserOrForbidden()
		reqCtx.SetResponseOk(true)
		reqCtx.SetResponseData(user)
	})
}
