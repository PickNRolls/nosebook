package presenterpost

import (
	"time"

	"github.com/google/uuid"
)

type postDest struct {
	Id        uuid.UUID `db:"id"`
	AuthorId  uuid.UUID `db:"author_id"`
	OwnerId   uuid.UUID `db:"owner_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *postDest) ID() uuid.UUID {
	return this.Id
}

func (this *postDest) Timestamp() time.Time {
	return this.CreatedAt
}
