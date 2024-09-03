package presenterfriendship

import (
	"time"

	"github.com/google/uuid"
)

type describe_relation_dest struct {
	RequesterId uuid.UUID `db:"requester_id"`
	ResponderId uuid.UUID `db:"responder_id"`
	CreatedAt   time.Time `db:"created_at"`
	Accepted    bool      `db:"accepted"`
}
