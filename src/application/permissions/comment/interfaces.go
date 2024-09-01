package permissionscomment

import "github.com/google/uuid"

type CommentToUpdate interface {
	AuthorId() uuid.UUID
}

type CommentToRemove interface {
	AuthorId() uuid.UUID
	ResourceOwnerId() uuid.UUID
}
