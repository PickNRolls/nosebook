package rootconvservice

import (
	"time"

	"github.com/google/uuid"
)

type chatDest struct {
	Id        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
