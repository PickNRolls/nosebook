package presenterlike

import "github.com/google/uuid"

type likeDest struct {
	ResourceId uuid.UUID `db:"resource_id"`
	UserId     uuid.UUID `db:"user_id"`
}
