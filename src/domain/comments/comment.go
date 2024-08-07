package comments

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID   `json:"id" db:"id"`
	AuthorId  uuid.UUID   `json:"authorId" db:"author_id"`
	Message   string      `json:"message" db:"message"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	RemovedAt time.Time   `json:"-" db:"removed_at"`
	PostId    uuid.UUID   `json:"-"`
	LikedBy   []uuid.UUID `json:"-"`

	Events []CommentEvent `json:"-"`
}

func NewComment(authorId uuid.UUID, message string) *Comment {
	return &Comment{
		Id:        uuid.New(),
		AuthorId:  authorId,
		Message:   message,
		CreatedAt: time.Now(),
		RemovedAt: time.Time{},
		PostId:    uuid.UUID{},
		LikedBy:   make([]uuid.UUID, 0),

		Events: make([]CommentEvent, 0),
	}
}

func (c *Comment) WithPost(id uuid.UUID) *Comment {
	c.PostId = id
	return c
}

func (c *Comment) Like(userId uuid.UUID) *Comment {
	for i, id := range c.LikedBy {
		if id == userId {
			c.LikedBy[i] = c.LikedBy[len(c.LikedBy)-1]
			c.LikedBy = c.LikedBy[:len(c.LikedBy)-1]
			c.Events = append(c.Events, NewCommentUnlikeEvent(userId))
			return c
		}
	}

	c.LikedBy = append(c.LikedBy, userId)
	c.Events = append(c.Events, NewCommentLikeEvent(userId))
	return c
}
