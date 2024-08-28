package middlewares

import (
	"net/http"
	"nosebook/src/errors"
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

		if errorsAny, exists := ctx.Get("errors"); exists {
			errs, ok := errorsAny.([]*errors.Error)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, responseDTO)
				return
			}

			responseDTO.Errors = errs
		}

		ctx.JSON(ctx.Writer.Status(), responseDTO)
	}
}
