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

type order struct{}

func (this *order) Column() string {
	return "created_at"
}
func (this *order) Timestamp(dest *Dest) time.Time {
	return dest.CreatedAt
}
func (this *order) Id(dest *Dest) uuid.UUID {
	return dest.Id
}
func (this *order) Asc() bool {
	return false
}
