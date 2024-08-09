package posts

import (
	"fmt"
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
				return
			}
		}

		authorId := ctx.Query("authorId")
		if authorId != "" {
			filter.AuthorId, err = uuid.Parse(authorId)
			if err != nil {
				ctx.Error(err)
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
		fmt.Println(result)
		if result.Err != nil {
			ctx.Error(result.Err)
			return
		}

		ctx.Set("data", result)
	}
}
