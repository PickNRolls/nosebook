package domaincomment

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CommentBuilder struct {
	id                uuid.UUID
	authorId          uuid.UUID
	message           string
	postId            uuid.UUID
	createdAt         time.Time
	removedAt         sql.NullTime
	permissions       permissions
	raiseCreatedEvent bool
}

func NewBuilder() *CommentBuilder {
	return &CommentBuilder{}
}

func (this *CommentBuilder) Build() *Comment {
	permissions := this.permissions
	if permissions == nil {
		permissions = &defaultPermissions{}
	}

	return newComment(
		this.id,
		this.authorId,
		this.message,
		this.postId,
		this.createdAt,
		this.removedAt,
		permissions,
		this.raiseCreatedEvent,
	)
}

func (this *CommentBuilder) Id(id uuid.UUID) *CommentBuilder {
	this.id = id
	return this
}

func (this *CommentBuilder) AuthorId(id uuid.UUID) *CommentBuilder {
	this.authorId = id
	return this
}

func (this *CommentBuilder) Message(message string) *CommentBuilder {
	this.message = message
	return this
}

func (this *CommentBuilder) PostId(id uuid.UUID) *CommentBuilder {
	this.postId = id
	return this
}

func (this *CommentBuilder) CreatedAt(t time.Time) *CommentBuilder {
	this.createdAt = t
	return this
}

func (this *CommentBuilder) RemovedAt(t time.Time) *CommentBuilder {
	this.removedAt.Valid = true
	this.removedAt.Time = t
	return this
}

func (this *CommentBuilder) Permissions(p permissions) *CommentBuilder {
	this.permissions = p
	return this
}

func (this *CommentBuilder) RaiseCreatedEvent() *CommentBuilder {
	this.raiseCreatedEvent = true
	return this
}
