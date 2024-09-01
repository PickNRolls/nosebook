package presenterpost

import (
	"time"

	"github.com/google/uuid"
)

type Dest struct {
	Id        uuid.UUID `db:"id"`
	AuthorId  uuid.UUID `db:"author_id"`
	OwnerId   uuid.UUID `db:"owner_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *Dest) ID() uuid.UUID {
	return this.Id
}

func (this *Dest) Timestamp() time.Time {
	return this.CreatedAt
}
