package comments

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id" db:"id"`
	AuthorId  uuid.UUID `json:"authorId" db:"author_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	RemovedAt time.Time `json:"removedAt" db:"removed_at"`
}

func NewComment(authorId uuid.UUID, message string) *Comment {
	return &Comment{
		Id:        uuid.New(),
		AuthorId:  authorId,
		Message:   message,
		CreatedAt: time.Now(),
		RemovedAt: time.Time{},
	}
}
