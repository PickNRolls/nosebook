package roothttp

import (
	userauth "nosebook/src/application/services/user_auth"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
)

func (this *RootHTTP) addAuthHandlers(userAuthService *userauth.Service) {
	this.unauthRouter.POST("/register", execResultHandler(userAuthService.RegisterUser))
	this.unauthRouter.POST("/login", execResultHandler(userAuthService.Login))

	this.authRouter.POST("/logout", execResultHandler(userAuthService.Logout, &execHandlerAvoidBinding{}))
	this.authRouter.GET("/whoami", func(ctx *gin.Context) {
		reqCtx := reqcontext.From(ctx)
		user := reqCtx.UserOrForbidden()
		reqCtx.SetResponseOk(true)
		reqCtx.SetResponseData(user)
	})
}
