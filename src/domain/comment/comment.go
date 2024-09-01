package domaincomment

import (
	"database/sql"
	"nosebook/src/lib/clock"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID
	AuthorId  uuid.UUID
	Message   string
	CreatedAt time.Time
	RemovedAt sql.NullTime
	PostId    uuid.UUID

	events []CommentEvent
}

func newComment(
	id uuid.UUID,
	authorId uuid.UUID,
	message string,
	postId uuid.UUID,
	createdAt time.Time,
	removedAt sql.NullTime,
	raiseCreatedEvent bool,
) *Comment {
	comment := &Comment{
		Id:        id,
		AuthorId:  authorId,
		Message:   message,
		PostId:    postId,
		CreatedAt: createdAt,
		RemovedAt: removedAt,

		events: make([]CommentEvent, 0),
	}

	if raiseCreatedEvent {
		comment.raiseEvent(NewCommentCreatedEvent())
	}

	return comment
}

func (this *Comment) raiseEvent(event CommentEvent) {
	this.events = append(this.events, event)
}

func (this *Comment) Events() []CommentEvent {
	return this.events
}

func (c *Comment) CanBeRemovedBy(userId uuid.UUID) *CommentError {
	if c.AuthorId != userId {
		return NewError("Только автор комментария может его удалить")
	}

	return nil
}

func (this *Comment) RemoveBy(userId uuid.UUID) *CommentError {
	err := this.CanBeRemovedBy(userId)
	if err != nil {
		return err
	}

	if this.RemovedAt.Valid {
		return NewError("Комментарий уже удален")
	}

	this.RemovedAt = sql.NullTime{
		Time:  clock.Now(),
		Valid: true,
	}

	this.raiseEvent(NewCommentRemovedEvent())
	return nil
}
