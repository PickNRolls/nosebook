package rootconvservice

import (
	"time"

	"github.com/google/uuid"
)

type chatDest struct {
	Id          uuid.UUID `db:"id"`
	LeftUserId  uuid.UUID `db:"left_user_id"`
	RightUserId uuid.UUID `db:"right_user_id"`
	CreatedAt   time.Time `db:"created_at"`
}
