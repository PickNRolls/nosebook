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

	permissions Permissions

	events []CommentEvent
}

func newComment(
	id uuid.UUID,
	authorId uuid.UUID,
	message string,
	postId uuid.UUID,
	createdAt time.Time,
	removedAt sql.NullTime,
	permissions Permissions,
	raiseCreatedEvent bool,
) *Comment {
	comment := &Comment{
		Id:        id,
		AuthorId:  authorId,
		Message:   message,
		PostId:    postId,
		CreatedAt: createdAt,
		RemovedAt: removedAt,

		permissions: permissions,

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

func (this *Comment) CanBeUpdatedBy(userId uuid.UUID) *Error {
	return this.permissions.CanUpdateBy(this, userId)
}

func (this *Comment) CanBeRemovedBy(userId uuid.UUID) *Error {
	return this.permissions.CanRemoveBy(this, userId)
}

func (this *Comment) RemoveBy(userId uuid.UUID) *Error {
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
