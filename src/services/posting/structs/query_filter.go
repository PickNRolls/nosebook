package structs

import (
	"time"

	"github.com/google/uuid"
)

type QueryFilter struct {
	OwnerId  uuid.UUID `db:"owner_id"`
	AuthorId uuid.UUID `db:"author_id"`
	Cursor   time.Time `db:"cursor"`
}
