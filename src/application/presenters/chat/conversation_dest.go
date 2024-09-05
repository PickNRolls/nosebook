package presenterchat

import (
	"time"

	"github.com/google/uuid"
)

type conv_dest struct {
	Id             uuid.UUID `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	LastMessageId  uuid.UUID `db:"last_message_id"`
	InterlocutorId uuid.UUID `db:"interlocutor_id"`
}

func (this *conv_dest) ID() uuid.UUID {
	return this.Id
}

func (this *conv_dest) Timestamp() time.Time {
	return this.CreatedAt
}
