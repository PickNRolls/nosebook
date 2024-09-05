package presentercomment

import (
	"time"

	"github.com/google/uuid"
)

type Dest struct {
	Id        uuid.UUID `db:"id"`
	PostId    uuid.UUID `db:"post_id"`
	AuthorId  uuid.UUID `db:"author_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
