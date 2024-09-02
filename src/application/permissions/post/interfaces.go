package permissionspost

import "github.com/google/uuid"

type Post interface {
	AuthorId() uuid.UUID
	OwnerId() uuid.UUID
}
