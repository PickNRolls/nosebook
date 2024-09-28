package presenteruser

import (
	"time"

	"github.com/google/uuid"
)

type dest struct {
	Id              uuid.UUID  `db:"id"`
	FirstName       string     `db:"first_name"`
	LastName        string     `db:"last_name"`
	Nick            string     `db:"nick"`
	LastActivityAt  time.Time  `db:"last_activity_at"`
	AvatarUrl       string     `db:"avatar_url"`
	AvatarUpdatedAt *time.Time `db:"avatar_updated_at"`
	CreatedAt       time.Time  `db:"created_at"`
}

type order struct{}

func (this *order) Column() string {
	return "created_at"
}
func (this *order) Timestamp(dest *dest) time.Time {
	return dest.CreatedAt
}
func (this *order) Id(dest *dest) uuid.UUID {
	return dest.Id
}
func (this *order) Asc() bool {
	return false
}
