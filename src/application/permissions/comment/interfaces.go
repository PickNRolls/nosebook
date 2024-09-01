package permissionscomment

import "github.com/google/uuid"

type Comment interface {
	AuthorId() uuid.UUID
}
