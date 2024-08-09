package middlewares

import (
	"net/http"
	"nosebook/src/infra/middlewares/presenter_middleware/dto"

	"github.com/gin-gonic/gin"
)

func NewPresenterMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		responseDTO := &dto.ResponseDTO{}

		data, exists := ctx.Get("data")
		if exists {
			responseDTO.Data = data
		}

		errorsAny, exists := ctx.Get("errors")
		if exists {
			errors, ok := errorsAny.([]error)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, responseDTO)
				return
			}

			responseDTO.Errors = errors
		}

		ctx.JSON(http.StatusOK, responseDTO)
	}
}
