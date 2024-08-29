package like

import "github.com/google/uuid"

type LikePostCommand struct {
	Id uuid.UUID `json:"id"`
}

type LikeCommentCommand struct {
	Id uuid.UUID `json:"id"`
}
