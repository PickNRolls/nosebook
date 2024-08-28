package comments

import (
	"database/sql"
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

func (c *Comment) WithPostId(id uuid.UUID) *Comment {
	c.PostId = id
	return c
}

func (c *Comment) Like(userId uuid.UUID) {
	for i, id := range c.LikedBy {
		if id == userId {
			c.LikedBy[i] = c.LikedBy[len(c.LikedBy)-1]
			c.LikedBy = c.LikedBy[:len(c.LikedBy)-1]
			c.Events = append(c.Events, NewCommentUnlikeEvent(userId))
		}
	}

	c.LikedBy = append(c.LikedBy, userId)
	c.Events = append(c.Events, NewCommentLikeEvent(userId))
}

func (c *Comment) CanBeRemovedBy(userId uuid.UUID) *CommentError {
	if c.AuthorId != userId {
		return NewError("Только автор комментария может его удалить")
	}

	return nil
}

func (c *Comment) Remove(userId uuid.UUID) *CommentError {
	err := c.CanBeRemovedBy(userId)
	if err != nil {
		return err
	}

	if c.RemovedAt.Valid {
		return NewError("Комментарий уже удален")
	}

	c.RemovedAt = sql.Null[time.Time]{
		V:     time.Now(),
		Valid: true,
	}
	c.Events = append(c.Events, NewCommentRemoveEvent())
	return nil
}
