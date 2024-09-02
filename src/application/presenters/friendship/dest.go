package presenterfriendship

import (
	"time"

	"github.com/google/uuid"
)

type dest struct {
	Id        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *dest) ID() uuid.UUID {
	return this.Id
}

func (this *dest) Timestamp() time.Time {
	return this.CreatedAt
}
