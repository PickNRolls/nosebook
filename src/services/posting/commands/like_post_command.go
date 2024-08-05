package commands

import "github.com/google/uuid"

type LikePostCommand struct {
	PostId uuid.UUID `json:"postId"`
}
