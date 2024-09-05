package presentercomment

import (
	"time"

	"github.com/google/uuid"
)

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
	return true
}
