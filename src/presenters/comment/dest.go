package presentercomment

import (
	"time"

	"github.com/google/uuid"
)

type dest struct {
	Id        uuid.UUID `db:"id"`
	PostId    uuid.UUID `db:"post_id"`
	AuthorId  uuid.UUID `db:"author_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *dest) ID() uuid.UUID {
	return this.Id
}

func (this *dest) Timestamp() time.Time {
	return this.CreatedAt
}
