package middlewares

import (
	"net/http"
	"nosebook/src/infra/errors"

	"github.com/gin-gonic/gin"
)

func NewNoRouteHandler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusNotFound)
		ctx.Error(errors.NewRouteNotFoundError())
	}
}
