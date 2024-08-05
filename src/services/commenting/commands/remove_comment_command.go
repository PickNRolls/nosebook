package commands

import (
	"github.com/google/uuid"
)

type RemoveCommentCommand struct {
	CommentId uuid.UUID `json:"commentId"`
}
