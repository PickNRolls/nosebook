package commenting

import (
	"github.com/google/uuid"
)

type PublishPostCommentCommand struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

type RemoveCommentCommand struct {
	Id uuid.UUID `json:"id"`
}
