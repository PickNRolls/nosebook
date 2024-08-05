package commands

import (
	"github.com/google/uuid"
)

type PublishPostCommentCommand struct {
	PostId  uuid.UUID `json:"postId"`
	Message string    `json:"message"`
}
