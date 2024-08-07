package dto

import (
	"github.com/google/uuid"
)

type QueryFilterDTO struct {
	OwnerId  uuid.UUID
	AuthorId uuid.UUID
	Cursor   string
}
