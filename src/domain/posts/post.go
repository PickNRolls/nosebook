package posts

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID `json:"id" db:"id"`
	AuthorId  uuid.UUID `json:"authorId" db:"author_id"`
	OwnerId   uuid.UUID `json:"ownerId" db:"owner_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	RemovedAt time.Time `json:"removedAt" db:"removed_at"`
}

func NewPost(authorId uuid.UUID, ownerId uuid.UUID, message string) *Post {
	return &Post{
		Id:        uuid.New(),
		AuthorId:  authorId,
		OwnerId:   ownerId,
		Message:   message,
		CreatedAt: time.Now(),
		RemovedAt: time.Time{},
	}
}
