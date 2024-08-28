//go:build exclude

package comments

import (
	"fmt"
	"nosebook/src/services"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerFind(commentService *services.CommentService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		filter := structs.QueryFilter{}
		var err error

		postIdsQuery := ctx.QueryArray("postIds")
		postIds := make([]uuid.UUID, len(postIdsQuery))
		fmt.Println(postIdsQuery)
		for i, p := range postIdsQuery {
			postIds[i], err = uuid.Parse(p)
			if err != nil {
				ctx.Error(err)
				return
			}
		}
		if len(postIds) > 0 {
			result := commentService.BatchFindByPostIds(postIds)
			ctx.Set("data", result)
			return
		}

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
		})
		if result.Err != nil {
			ctx.Error(result.Err)
			return
		}

		ctx.Set("data", result)
	}
}
