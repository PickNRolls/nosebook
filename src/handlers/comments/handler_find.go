package comments

import (
	"nosebook/src/infra/helpers"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerFind(commentService *services.CommentService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrForbidden(ctx)

		filter := structs.QueryFilter{}
		var err error

		postId := ctx.Query("postId")
		if postId != "" {
			filter.PostId, err = uuid.Parse(postId)
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

		next := ctx.Query("next")
		if next != "" {
			filter.Next = next
		}

		prev := ctx.Query("prev")
		if prev != "" {
			filter.Prev = prev
		}

		last := ctx.Query("last")
		if last != "" {
			filter.Last = true
		}

		result := commentService.FindByFilter(&commands.FindCommentsCommand{
			Filter: filter,
		}, &auth.Auth{
			UserId: user.ID,
		})
		if result.Err != nil {
			ctx.Error(result.Err)
			return
		}

		ctx.Set("data", result)
	}
}
