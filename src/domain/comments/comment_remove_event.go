package comments

import "github.com/google/uuid"

type CommentRemoveEvent struct {
	UserId uuid.UUID
}

func (event *CommentRemoveEvent) Type() CommentEventType {
	return REMOVED
}

func NewCommentRemoveEvent() *CommentRemoveEvent {
	return &CommentRemoveEvent{}
}
