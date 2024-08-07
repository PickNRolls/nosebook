package posts

import (
	"net/http"
	"nosebook/src/infra/helpers"
	"nosebook/src/services"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewHandlerFind(postingService *services.PostingService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user := helpers.GetUserOrBadRequest(ctx)

		command := commands.FindPostsCommand{}
		var err error

		ownerId := ctx.Query("ownerId")
		if ownerId != "" {
			command.Filter.OwnerId, err = uuid.Parse(ownerId)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		authorId := ctx.Query("authorId")
		if authorId != "" {
			command.Filter.AuthorId, err = uuid.Parse(authorId)
			if err != nil {
				ctx.Error(err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		result := postingService.FindByFilter(&command, &auth.Auth{
			UserId: user.ID,
		})
		if result.Err != nil {
			ctx.Error(result.Err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
