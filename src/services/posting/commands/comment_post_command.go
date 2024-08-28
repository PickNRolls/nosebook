package commands

import "github.com/google/uuid"

type CommentPostCommand struct {
	Id      uuid.UUID `json:"id"`
	Comment string    `json:"message"`
}
