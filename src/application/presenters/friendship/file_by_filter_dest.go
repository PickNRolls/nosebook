package presenterfriendship

import (
	"time"

	"github.com/google/uuid"
)

type find_by_filter_dest struct {
	Id        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

func (this *find_by_filter_dest) ID() uuid.UUID {
	return this.Id
}

func (this *find_by_filter_dest) Timestamp() time.Time {
	return this.CreatedAt
}
