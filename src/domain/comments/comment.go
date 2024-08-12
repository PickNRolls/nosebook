package comments

import (
	"database/sql"
	"nosebook/src/errors"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID           `json:"id" db:"id"`
	AuthorId  uuid.UUID           `json:"authorId" db:"author_id"`
	Message   string              `json:"message" db:"message"`
	CreatedAt time.Time           `json:"createdAt" db:"created_at"`
	RemovedAt sql.Null[time.Time] `json:"-" db:"removed_at"`
	PostId    uuid.UUID           `json:"-" db:"post_id"`
	LikedBy   []uuid.UUID         `json:"-"`

	Events []CommentEvent `json:"-"`
}

func NewComment(authorId uuid.UUID, message string) *Comment {
	return &Comment{
		Id:        uuid.New(),
		AuthorId:  authorId,
		Message:   message,
		CreatedAt: time.Now(),
		RemovedAt: sql.Null[time.Time]{},
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

func (c *Comment) Remove(userId uuid.UUID) (*Comment, *errors.Error) {
	if c.AuthorId != userId {
		return nil, errors.New("CommentError", "Только автор комментария может его удалить")
	}

	if c.RemovedAt.Valid {
		return nil, errors.New("CommentError", "Комментарий уже удален")
	}

	c.RemovedAt = sql.Null[time.Time]{
		V:     time.Now(),
		Valid: true,
	}
	c.Events = append(c.Events, NewCommentRemoveEvent())
	return c, nil
}
