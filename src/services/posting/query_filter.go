package posting

import (
	"github.com/google/uuid"
)

type QueryFilter struct {
	OwnerId  uuid.UUID
	AuthorId uuid.UUID
	Cursor   string
}
