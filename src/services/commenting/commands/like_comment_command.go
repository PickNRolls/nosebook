package commands

import "github.com/google/uuid"

type LikeCommentCommand struct {
	CommentId uuid.UUID `json:"commentId"`
}
