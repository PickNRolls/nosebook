package commands

import "github.com/google/uuid"

type LikePostCommand struct {
	Id uuid.UUID `json:"id"`
}
