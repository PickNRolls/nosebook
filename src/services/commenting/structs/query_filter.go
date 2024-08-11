package structs

import (
	"github.com/google/uuid"
)

type QueryFilter struct {
	PostId   uuid.UUID
	AuthorId uuid.UUID
	Next     string
	Prev     string
	Last     bool
}
