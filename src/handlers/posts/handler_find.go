package posts

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/presenters"
	"nosebook/src/presenters/post_presenter/dto"
	"nosebook/src/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerFind(postPresenter *presenters.PostPresenter) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)

		filter := dto.QueryFilterDTO{}
		var err error

		ownerId := ctx.Query("ownerId")
		if ownerId != "" {
			filter.OwnerId, err = uuid.Parse(ownerId)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		authorId := ctx.Query("authorId")
		if authorId != "" {
			filter.AuthorId, err = uuid.Parse(authorId)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		cursor := ctx.Query("cursor")
		if cursor != "" {
			filter.Cursor = cursor
		}

		result := postPresenter.FindByFilter(filter, &auth.Auth{
			UserId: user.ID,
		})
		if result.Err != nil {
			ctx.Error(result.Err)
			ctx.JSON(http.StatusBadRequest, result)
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
