package domainfriendship

import (
	"time"

	"github.com/google/uuid"
)

type Builder struct {
	requesterId       uuid.UUID
	responderId       uuid.UUID
	message           string
	accepted          bool
	viewed            bool
	createdAt         time.Time
	raiseCreatedEvent bool
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (this *Builder) RequesterId(id uuid.UUID) *Builder {
	this.requesterId = id
	return this
}

func (this *Builder) ResponderId(id uuid.UUID) *Builder {
	this.responderId = id
	return this
}

func (this *Builder) Message(message string) *Builder {
	this.message = message
	return this
}

func (this *Builder) Accepted(accepted bool) *Builder {
	this.accepted = accepted
	return this
}

func (this *Builder) Viewed(viewed bool) *Builder {
	this.viewed = viewed
	return this
}

func (this *Builder) CreatedAt(createdAt time.Time) *Builder {
	this.createdAt = createdAt
	return this
}

func (this *Builder) RaiseCreatedEvent() *Builder {
	this.raiseCreatedEvent = true
	return this
}

func (this *Builder) Build() *FriendRequest {
	return New(
		this.requesterId,
		this.responderId,
		this.message,
		this.accepted,
		this.viewed,
		this.createdAt,
		this.raiseCreatedEvent,
	)
}
