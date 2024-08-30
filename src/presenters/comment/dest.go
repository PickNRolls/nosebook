package presentercomment

import (
	"time"

	"github.com/google/uuid"
)

type commentDest struct {
	Id        uuid.UUID `db:"id"`
	PostId    uuid.UUID `db:"post_id"`
	AuthorId  uuid.UUID `db:"author_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *commentDest) ID() uuid.UUID {
	return this.Id
}

func (this *commentDest) Timestamp() time.Time {
	return this.CreatedAt
}
