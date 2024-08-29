package domaincomment

import "github.com/google/uuid"

type CommentRemovedEvent struct {
	RemoverUserId uuid.UUID
}

func (event *CommentRemovedEvent) Type() CommentEventType {
	return REMOVED
}

func NewCommentRemovedEvent() *CommentRemovedEvent {
	return &CommentRemovedEvent{}
}
