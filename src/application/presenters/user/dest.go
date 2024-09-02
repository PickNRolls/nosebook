package presenteruser

import (
	"time"

	"github.com/google/uuid"
)

type dest struct {
	Id             uuid.UUID `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	Nick           string    `db:"nick"`
	LastActivityAt time.Time `db:"last_activity_at"`
}
