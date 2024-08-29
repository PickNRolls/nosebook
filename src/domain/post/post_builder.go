package domainpost

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PostBuilder struct {
	id                uuid.UUID
	authorId          uuid.UUID
	ownerId           uuid.UUID
	message           string
	createdAt         time.Time
	removedAt         sql.NullTime
	raiseCreatedEvent bool
}

func NewBuilder() *PostBuilder {
	return &PostBuilder{}
}

func (this *PostBuilder) Build() *Post {
	return NewPost(
		this.id,
		this.authorId,
		this.ownerId,
		this.message,
		this.createdAt,
		this.removedAt,
		this.raiseCreatedEvent,
	)
}

func (this *PostBuilder) Id(id uuid.UUID) *PostBuilder {
	this.id = id
	return this
}

func (this *PostBuilder) AuthorId(id uuid.UUID) *PostBuilder {
	this.authorId = id
	return this
}

func (this *PostBuilder) OwnerId(id uuid.UUID) *PostBuilder {
	this.ownerId = id
	return this
}

func (this *PostBuilder) Message(message string) *PostBuilder {
	this.message = message
	return this
}

func (this *PostBuilder) CreatedAt(t time.Time) *PostBuilder {
	this.createdAt = t
	return this
}

func (this *PostBuilder) RemovedAt(t time.Time) *PostBuilder {
	this.removedAt.Valid = true
	this.removedAt.Time = t
	return this
}

func (this *PostBuilder) RaiseCreatedEvent() *PostBuilder {
	this.raiseCreatedEvent = true
	return this
}
