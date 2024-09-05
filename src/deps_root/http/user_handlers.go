package roothttp

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/presenters/user"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (this *RootHTTP) addUserHandlers() {
	presenter := presenteruser.New(this.db)

	group := this.authRouter.Group("/users")

	group.GET("/:id", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)
		id, ok := reqctx.ParamUUID("id")
		if !ok {
			return
		}

		users, ok := handle(presenter.FindByIds([]uuid.UUID{id}))(reqctx)
		if !ok {
			return
		}

		var user *presenterdto.User
		if len(users) > 0 {
			user = users[id]
		}

		reqctx.SetResponseData(user)
		reqctx.SetResponseOk(true)
	})
}
