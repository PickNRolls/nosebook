package dto

import (
	"github.com/google/uuid"
)

type FindInputDTO struct {
	OwnerId  uuid.UUID
	AuthorId uuid.UUID
	Cursor   string
}
