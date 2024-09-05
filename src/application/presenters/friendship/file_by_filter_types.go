package presenterfriendship

import (
	"time"

	"github.com/google/uuid"
)

type find_by_filter_dest struct {
	Id        uuid.UUID   `db:"id"`
	Type      RequestType `db:"type"`
	CreatedAt time.Time   `db:"created_at"`
	Accepted  bool        `db:"accepted"`
}

type find_by_filter_order struct{}

func (this *find_by_filter_order) Column() string {
	return "created_at"
}
func (this *find_by_filter_order) Timestamp(dest *find_by_filter_dest) time.Time {
	return dest.CreatedAt
}
func (this *find_by_filter_order) Id(dest *find_by_filter_dest) uuid.UUID {
	return dest.Id
}
func (this *find_by_filter_order) Asc() bool {
	return false
}
