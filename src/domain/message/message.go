package domainmessage

import (
	"database/sql"
	"nosebook/src/lib/clock"
	"nosebook/src/lib/nullable"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID
	AuthorId  uuid.UUID
	Text      string
	ReplyTo   uuid.UUID
	CreatedAt time.Time
	RemovedAt nullable.Time

	permissions Permissions
	events      []Event
}

func New(
	id uuid.UUID,
	authorId uuid.UUID,
	text string,
	replyTo uuid.UUID,
	createdAt time.Time,
	removedAt sql.NullTime,
	permissions Permissions,
	raiseCreatedEvent bool,
) *Message {
	message := &Message{
		Id:        id,
		AuthorId:  authorId,
		Text:      text,
		ReplyTo:   replyTo,
		CreatedAt: createdAt,
		RemovedAt: removedAt,

		permissions: permissions,

		events: make([]Event, 0),
	}

	if message.permissions == nil {
		message.permissions = &defaultPermissions{}
	}

	if raiseCreatedEvent {
		message.raiseEvent(NewCreatedEvent())
	}

	return message
}

func (this *Message) raiseEvent(event Event) {
	this.events = append(this.events, event)
}

func (this *Message) Events() []Event {
	return this.events
}

func (this *Message) CanBeUpdatedBy(userId uuid.UUID) *Error {
	return this.permissions.CanUpdateBy(this, userId)
}

func (this *Message) CanBeRemovedBy(userId uuid.UUID) *Error {
	return this.permissions.CanRemoveBy(this, userId)
}

func (this *Message) RemoveBy(userId uuid.UUID) *Error {
	err := this.CanBeRemovedBy(userId)
	if err != nil {
		return err
	}

	if this.RemovedAt.Valid {
		return newError("Сообщение уже удалено")
	}

	this.RemovedAt = nullable.Time{
		Time:  clock.Now(),
		Valid: true,
	}

	this.raiseEvent(NewRemovedEvent())
	return nil
}
