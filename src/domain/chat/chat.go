package domainchat

import (
	domainmessage "nosebook/src/domain/message"
	"nosebook/src/lib/clock"
	"nosebook/src/lib/nullable"

	"github.com/google/uuid"
)

type Chat struct {
	Id          uuid.UUID
	MemberIds   uuid.UUIDs
	Name        string
	Private     bool
	permissions Permissions
	events      []Event
}

func New(
	id uuid.UUID,
	memberIds uuid.UUIDs,
	name string,
	private bool,
	permissions Permissions,
	raiseCreatedEvent bool,
) (*Chat, *Error) {
	if private && len(memberIds) > 2 {
		return nil, newError("Приватный чат не может содержать больше 2 участников")
	}

	out := &Chat{
		Id:          id,
		MemberIds:   memberIds,
		Name:        name,
		Private:     private,
		permissions: permissions,
	}

	if out.permissions == nil {
		out.permissions = &defaultPermissions{}
	}

	if raiseCreatedEvent {
		out.raiseEvent(newCreatedEvent())
	}

	return out, nil
}

func (this *Chat) Events() []Event {
	return this.events
}

func (this *Chat) raiseEvent(event Event) {
	this.events = append(this.events, event)
}

func (this *Chat) CanJoin(userId uuid.UUID) *Error {
	return this.permissions.CanJoinBy(this, userId)
}

func (this *Chat) CanSendMessageBy(userId uuid.UUID) *Error {
	return this.permissions.CanSendMessageBy(this, userId)
}

func (this *Chat) SendMessageBy(text string, replyTo uuid.UUID, userId uuid.UUID) *Error {
	if err := this.CanSendMessageBy(userId); err != nil {
		return err
	}

	message := domainmessage.New(
		uuid.New(),
		userId,
		text,
		replyTo,
		clock.Now(),
		nullable.Time{},
		nil,
		true,
	)
	this.raiseEvent(newMessageSentEvent(message))

	return nil
}
